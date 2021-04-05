package httpext

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
)

type tlsConfigBuilder struct {
	insecureSkipVerify bool
	caCertFiles        []string
	certFile           string
	keyFile            string
}

// TLSConfigBuilder - interface for the builder of http.Transport
type TLSConfigBuilder interface {
	// SkipVerify - set the InsecureSkipVerify flag
	SkipVerify(insecureSkipVerify bool) TLSConfigBuilder
	// CACerts - set the CA Cert files
	CACerts(caCertFiles ...string) TLSConfigBuilder
	// KeyPair - Set the X509 public (cert) and private (key) files
	KeyPair(certFile string, keyFile string) TLSConfigBuilder
	// Build - build the tls.Config
	Build() (*tls.Config, error)
}

// NewTLSConfigBuilder - create a new HTTPTransportBuilder
func NewTLSConfigBuilder() TLSConfigBuilder {
	return &tlsConfigBuilder{}
}

func (tb *tlsConfigBuilder) SkipVerify(insecureSkipVerify bool) TLSConfigBuilder {
	tb.insecureSkipVerify = insecureSkipVerify
	return tb
}

func (tb *tlsConfigBuilder) CACerts(caCertFiles ...string) TLSConfigBuilder {
	tb.caCertFiles = caCertFiles
	return tb
}

func (tb *tlsConfigBuilder) KeyPair(certFile string, keyFile string) TLSConfigBuilder {
	tb.certFile = certFile
	tb.keyFile = keyFile
	return tb
}

func checkFilesExist(fileType string, files ...string) error {
	if len(files) > 0 {
		for _, file := range files {
			if file != "" && !fileExists(file) {
				return fmt.Errorf("%s file '%s' does not exist", fileType, file)
			}
		}
	}

	return nil
}

func (tb *tlsConfigBuilder) Build() (*tls.Config, error) {
	// Check that files exist
	if err := checkFilesExist("Cert", tb.certFile); err != nil {
		return nil, err
	}
	if err := checkFilesExist("Key", tb.keyFile); err != nil {
		return nil, err
	}
	if err := checkFilesExist("CA", tb.caCertFiles...); err != nil {
		return nil, err
	}

	var tlsCfg tls.Config

	tlsCfg.InsecureSkipVerify = tb.insecureSkipVerify

	if tb.certFile != "" && tb.keyFile != "" {
		// Load the cert,key and ca files
		clientCert, err := tls.LoadX509KeyPair(tb.certFile, tb.keyFile)
		if err != nil {
			return nil, err
		}

		// At this point, we have loaded the cert,key and ca files
		tlsCfg.Certificates = append(tlsCfg.Certificates, clientCert)
	}

	if len(tb.caCertFiles) > 0 {
		caCertPool := x509.NewCertPool()

		for _, caCertFile := range tb.caCertFiles {
			caCertBytes, err := ioutil.ReadFile(caCertFile)
			if err != nil {
				return nil, err
			}
			if ok := caCertPool.AppendCertsFromPEM(caCertBytes); !ok {
				return nil, fmt.Errorf("failed to append certs from: %s", caCertFile)
			}
		}
		tlsCfg.RootCAs = caCertPool
	}

	return &tlsCfg, nil
}

// fileExists - check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

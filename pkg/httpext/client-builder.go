package httpext

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

// "github.com/sanity-io/litter"

type CheckRetryFunc retryablehttp.CheckRetry
type BackoffFunc retryablehttp.Backoff

// type CheckRetryFunc func(resp *http.Response, err error) (bool, error)
// type BackoffFunc func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration

// var mylitter = litter.Options{
// 	HidePrivateFields: false,
// 	// HomePackage: "thispack",
// 	// Separator: " ",
// }

type httpClientBuilder struct {
	Foo        int
	tlsBuilder TLSConfigBuilder
	tlsCfg     *tls.Config

	retries                    int
	retryMinWait, retryMaxWait time.Duration
	checkRetryFn               CheckRetryFunc
	backoffFn                  BackoffFunc
}

type HttpClientBuilder interface {
	Retries(retries int) HttpClientBuilder
	CheckRetry(checkRetryFn CheckRetryFunc) HttpClientBuilder
	Backoff(backoffFn BackoffFunc) HttpClientBuilder

	// TLSConfig - with a specific TLS configuration
	TLSConfig(tlsCfg *tls.Config) HttpClientBuilder
	// TLSSkipVerify - set the InsecureSkipVerify flag
	TLSSkipVerify(insecureSkipVerify bool) HttpClientBuilder
	// TLSCACerts - set the CA Cert files
	TLSCACerts(caCertFiles ...string) HttpClientBuilder
	// TLSKeyPair - Set the X509 public (cert) and private (key) files
	TLSKeyPair(certFile string, keyFile string) HttpClientBuilder

	Build() (*http.Client, error)
}

func NewHttpClientBuilder() HttpClientBuilder {
	return &httpClientBuilder{}
}

func (b *httpClientBuilder) Retries(retries int) HttpClientBuilder {
	if retries >= 0 {
		b.retries = retries
	}
	// mylitter.Dump("b.Retries()", b)
	return b
}

func (b *httpClientBuilder) CheckRetry(checkRetryFn CheckRetryFunc) HttpClientBuilder {
	b.checkRetryFn = checkRetryFn
	return b
}

func (b *httpClientBuilder) Backoff(backoffFn BackoffFunc) HttpClientBuilder {
	b.backoffFn = backoffFn
	return b
}

func (b *httpClientBuilder) TLSConfig(tlsCfg *tls.Config) HttpClientBuilder {
	b.tlsCfg = tlsCfg
	return b
}

func (b *httpClientBuilder) getTLSConfigBuilder() TLSConfigBuilder {
	if b.tlsBuilder == nil {
		b.tlsBuilder = NewTLSConfigBuilder()
	}
	return b.tlsBuilder
}
func (b *httpClientBuilder) TLSSkipVerify(insecureSkipVerify bool) HttpClientBuilder {
	b.getTLSConfigBuilder().SkipVerify(insecureSkipVerify)
	return b
}
func (b *httpClientBuilder) TLSCACerts(caCertFiles ...string) HttpClientBuilder {
	b.getTLSConfigBuilder().CACerts(caCertFiles...)
	return b
}

func (b *httpClientBuilder) TLSKeyPair(certFile string, keyFile string) HttpClientBuilder {
	b.getTLSConfigBuilder().KeyPair(certFile, keyFile)
	return b
}

func (b *httpClientBuilder) Build() (client *http.Client, err error) {
	if client, err = b.buildDefaultClient(); err != nil {
		return nil, err
	}

	// mylitter.Dump(b)
	var retryClient *retryablehttp.Client
	if false {
		retryClient = retryablehttp.NewClient()
		retryClient.HTTPClient = client

	} else {
		retryClient = &retryablehttp.Client{
			HTTPClient:   client,
			RetryMax:     3,
			RetryWaitMin: 1 * time.Second,
			RetryWaitMax: 30 * time.Second,
			CheckRetry:   retryablehttp.DefaultRetryPolicy,
			Backoff:      retryablehttp.DefaultBackoff,
			Logger:       log.New(os.Stderr, "", log.LstdFlags),
		}
	}

	if b.retries > 0 {
		retryClient.RetryMax = b.retries
	}
	if b.retryMaxWait > 0 {
		retryClient.RetryWaitMin = b.retryMaxWait
	}
	if b.retryMinWait > 0 {
		retryClient.RetryWaitMax = b.retryMinWait
	}
	if b.checkRetryFn != nil {
		retryClient.CheckRetry = retryablehttp.CheckRetry(b.checkRetryFn)
		// retryClient.CheckRetry = func(_ context.Context, resp *http.Response, err error) (bool, error) {
		// 	return b.checkRetryFn(resp, err)
		// }
	}
	if b.backoffFn != nil {
		retryClient.Backoff = retryablehttp.Backoff(b.backoffFn)
		// retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		// 	return b.backoffFn(min, max, attemptNum, resp)
		// }
	}

	// Build the HTTP client to return (replace roundtripper with one with retries)
	client = &http.Client{
		Transport: &retryablehttp.RoundTripper{
			Client: retryClient,
		},
	}

	// litter.Dump(retryClient)
	return
}

func (b *httpClientBuilder) buildDefaultClient() (client *http.Client, err error) {
	// Initialize client
	client = &http.Client{}

	// Build TLS Config
	var tlsCfg *tls.Config

	if b.tlsBuilder != nil {
		if tlsCfg, err = b.tlsBuilder.Build(); err != nil {
			return nil, errors.Wrap(err, "Failed to create TLS config")
		}
	}
	if tlsCfg == nil && b.tlsCfg != nil {
		tlsCfg = b.tlsCfg
	}

	if tlsCfg != nil {
		client.Transport = &http.Transport{
			TLSClientConfig: tlsCfg,
		}
	}

	return client, err
}

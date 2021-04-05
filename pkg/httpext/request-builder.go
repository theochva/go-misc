package httpext

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type requestBuilder struct {
	url         string
	method      string
	headers     map[string]string
	queryParams url.Values
	body        interface{}
}

// RequestBuilder - create a request builder
type RequestBuilder interface {
	Method(method string) RequestBuilder
	URL(url string) RequestBuilder
	Query(key string, value string) RequestBuilder
	Header(key, value string) RequestBuilder
	Headers(headers map[string]string) RequestBuilder
	Body(body interface{}) RequestBuilder
	Build() (*http.Request, error)
	BuildWithCtx(ctx context.Context) (*http.Request, error)
}

// NewRequestBuilder - build a request builder
func NewRequestBuilder() RequestBuilder {
	return &requestBuilder{headers: map[string]string{}, queryParams: url.Values{}}
}

func (rb *requestBuilder) Method(method string) RequestBuilder {
	rb.method = method
	return rb
}
func (rb *requestBuilder) Query(key string, value string) RequestBuilder {
	rb.queryParams.Add(key, value)
	return rb
}
func (rb *requestBuilder) URL(url string) RequestBuilder {
	rb.url = url
	return rb
}
func (rb *requestBuilder) Header(key, value string) RequestBuilder {
	rb.headers[key] = value
	return rb
}
func (rb *requestBuilder) Headers(headers map[string]string) RequestBuilder {
	rb.headers = map[string]string{}
	if len(headers) > 0 {
		for key, value := range headers {
			rb.headers[key] = value
		}
	}
	return rb
}

func (rb *requestBuilder) Body(body interface{}) RequestBuilder {
	rb.body = body
	return rb
}

func (rb *requestBuilder) Build() (req *http.Request, err error) {
	var (
		bodyStream io.Reader
		url        string
	)
	if rb.body != nil {
		if bodyStream, err = getBodyReader(rb.body); err != nil {
			return nil, err
		}
	}
	if url, err = rb.buildURL(); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest(rb.method, url, bodyStream); err != nil {
		return nil, err
	}

	if len(rb.headers) > 0 {
		for key, value := range rb.headers {
			req.Header.Add(key, value)
		}
	}
	return req, nil
}

func (rb *requestBuilder) BuildWithCtx(ctx context.Context) (req *http.Request, err error) {
	if req, err = rb.Build(); err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return
}

func (rb *requestBuilder) buildURL() (string, error) {
	if rb.url == "" || len(rb.queryParams) == 0 {
		return rb.url, nil
	}
	u, err := url.Parse(rb.url)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, vs := range rb.queryParams {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func getBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	switch bodyType := body.(type) {
	case string:
		return strings.NewReader(body.(string)), nil
	case []byte:
		return bytes.NewReader(body.([]byte)), nil
	case io.Reader:
		return body.(io.Reader), nil
	default:
		log.Debug("Body type: %v", bodyType)
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON request: %v", err)
		}
		return bytes.NewReader(raw), nil
	}
}

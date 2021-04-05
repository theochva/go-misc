package httpext

import "net/http"

// NewGetRequestGuilder - build a request builder with method set to GET
func NewGetRequestBuilder() RequestBuilder {
	return NewRequestBuilder().Method(http.MethodGet)
}

// NewPostRequestBuilder - build a request builder with method set to POST
func NewPostRequestBuilder() RequestBuilder {
	return NewRequestBuilder().Method(http.MethodPost)
}

// NewPatchRequestBuilder - build a request builder with method set to PATCH
func NewPatchRequestBuilder() RequestBuilder {
	return NewRequestBuilder().Method(http.MethodPatch)
}

// NewPutRequestBuilder - build a request builder with method set to PUT
func NewPutRequestBuilder() RequestBuilder {
	return NewRequestBuilder().Method(http.MethodPut)
}

// NewDeleteRequestBuilder - build a request builder with method set to DELETE
func NewDeleteRequestBuilder() RequestBuilder {
	return NewRequestBuilder().Method(http.MethodDelete)
}

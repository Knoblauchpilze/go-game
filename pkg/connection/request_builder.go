package connection

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type RequestBuilder struct {
	url     string
	headers http.Header
	builder httpRequestBuilder
	body    interface{}
	client  HttpClient
}

func newRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		headers: make(http.Header),
		client:  &http.Client{},
	}
}

func (rb *RequestBuilder) SetUrl(url string) {
	rb.url = url
}

func (rb *RequestBuilder) SetHeaders(headers http.Header) {
	rb.headers = headers
}

func (rb *RequestBuilder) AddHeader(key string, value []string) {
	existing := rb.headers[key]
	existing = append(existing, value...)
	rb.headers[key] = existing
}

func (rb *RequestBuilder) SetBody(contentType string, body interface{}) {
	rb.body = body
	rb.AddHeader("Content-Type", []string{contentType})
}

func (rb *RequestBuilder) setHttpRequestBuilder(builder httpRequestBuilder) {
	rb.builder = builder
}

func (rb *RequestBuilder) setHttpClient(client HttpClient) {
	rb.client = client
}

func (rb *RequestBuilder) Build() (RequestWrapper, error) {
	if len(rb.url) == 0 {
		return nil, errors.New("invalid url for request")
	}
	if rb.builder == nil {
		return nil, errors.New("invalid builder for request")
	}

	ri := &requestImpl{
		url:     rb.url,
		headers: rb.headers,
		builder: rb.builder,
		body:    rb.body,
		client:  rb.client,
	}

	return ri, nil
}

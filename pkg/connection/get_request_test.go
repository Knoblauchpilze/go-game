package connection

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpGetRequest_UrlReached(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{}
	url := "http://dummy-url"

	rb := NewHttpGetRequestBuilder()
	rb.SetUrl(url)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	rw.Perform()
	assert.Equal(url, mc.inReq.URL.String())
}

func TestHttpGetRequest_InvalidUrl(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{}

	rb := NewHttpGetRequestBuilder()
	rb.SetUrl(invalidUrl)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	_, err = rw.Perform()
	assert.NotNil(err)
	assert.Nil(mc.inReq)
}

func TestHttpGetRequest_HeadersPassed(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{}
	url := "http://dummy-url"
	headers := http.Header{
		"haha": []string{"jaja"},
	}

	rb := NewHttpGetRequestBuilder()
	rb.SetUrl(url)
	rb.SetHeaders(headers)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	rw.Perform()
	assert.Equal(headers, mc.inReq.Header)
}

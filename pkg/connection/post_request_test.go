package connection

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpPostRequest_UrlReached(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{}
	url := "http://dummy-url"

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl(url)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	rw.Perform()
	assert.Equal(url, mc.inReq.URL.String())
}

func TestHttpPostRequest_InvalidUrl(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{}

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl(invalidUrl)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	_, err = rw.Perform()
	assert.NotNil(err)
	assert.Nil(mc.inReq)
}

func TestHttpPostRequest_HeadersPassed(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{}
	url := "http://dummy-url"
	headers := http.Header{
		"haha": []string{"jaja"},
	}

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl(url)
	rb.SetHeaders(headers)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	rw.Perform()
	assert.Equal(headers, mc.inReq.Header)
}

func TestHttpPostRequest_BodyPassed(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{}
	url := "http://dummy-url"

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl(url)
	rb.SetBody("kiki", "some data")
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	rw.Perform()
	out, err := io.ReadAll(mc.inReq.Body)
	assert.Nil(err)
	assert.Equal("\"some data\"", string(out))
}

func TestHttpPostRequest_UnmarshallableBody(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{}
	url := "http://dummy-url"

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl(url)
	rb.SetBody("kiki", unmarshallableContent{})
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	_, err = rw.Perform()
	assert.NotNil(err)
}

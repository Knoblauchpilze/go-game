package connection

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_Perform_BuilderError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{
		expectedError: nil,
		expectedResp:  nil,
	}

	rb := newRequestBuilder()
	rb.SetUrl("haha")
	rb.setHttpClient(mc)
	rb.setHttpRequestBuilder(errorHttpBuilder)
	rw, err := rb.Build()
	assert.Nil(err)

	_, err = rw.Perform()
	assert.NotNil(err)
}

func TestRequest_Perform_BuilderInvalid(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{
		expectedError: nil,
		expectedResp:  nil,
	}

	rb := newRequestBuilder()
	rb.SetUrl("haha")
	rb.setHttpClient(mc)
	rb.setHttpRequestBuilder(nilRequestHttpBuilder)
	rw, err := rb.Build()
	assert.Nil(err)

	_, err = rw.Perform()
	assert.NotNil(err)
}

func TestRequest_Perform_ClientErrorReceived(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{
		expectedError: errSomeError,
		expectedResp:  nil,
	}

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl("haha")
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	_, err = rw.Perform()
	assert.Equal(mc.expectedError, err)
}

func TestRequest_Perform_ClientResponseReceived(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{
		expectedError: nil,
		expectedResp:  nil,
	}

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl("haha")
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	resp, _ := rw.Perform()
	assert.Equal(mc.expectedResp, resp)
}

func TestRequest_Perform_StatusCodeReceived(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{
		expectedError: nil,
		expectedResp:  generateHttpResponse(),
	}

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl("haha")
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	resp, _ := rw.Perform()
	assert.Equal(mc.expectedResp.StatusCode, resp.StatusCode)
}

func TestRequest_Perform_HeaderReceived(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{
		expectedError: nil,
		expectedResp:  generateHttpResponse(),
	}

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl("haha")
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	resp, _ := rw.Perform()
	assert.Equal(mc.expectedResp.Header, resp.Header)
}

func TestRequest_Perform_BodyReceived(t *testing.T) {
	assert := assert.New(t)

	mc := &mockHttpClient{
		expectedError: nil,
		expectedResp:  generateHttpResponse(),
	}

	rb := NewHttpPostRequestBuilder()
	rb.SetUrl("haha")
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	resp, _ := rw.Perform()
	out, err := io.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("some data", string(out))
}

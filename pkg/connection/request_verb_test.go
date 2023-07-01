package connection

import (
	"io"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type requestGenerator func() *RequestBuilder

// https://pkg.go.dev/github.com/stretchr/testify/suite
type RequestWithVerbTestSuite struct {
	suite.Suite
	verbAcceptsBody bool
	createRequest   requestGenerator
}

func (suite *RequestWithVerbTestSuite) TestHttpRequest_UrlReached() {
	assert := assert.New(suite.T())

	mc := &mockHttpClient{}
	url := "http://dummy-url"

	rb := suite.createRequest()
	rb.SetUrl(url)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	rw.Perform()
	assert.Equal(url, mc.inReq.URL.String())
}

func (suite *RequestWithVerbTestSuite) TestHttpRequest_InvalidUrl() {
	assert := assert.New(suite.T())

	mc := &mockHttpClient{}

	rb := suite.createRequest()
	rb.SetUrl(invalidUrl)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	_, err = rw.Perform()
	assert.NotNil(err)
	assert.Nil(mc.inReq)
}

func (suite *RequestWithVerbTestSuite) TestHttpRequest_HeadersPassed() {
	assert := assert.New(suite.T())

	mc := &mockHttpClient{}
	url := "http://dummy-url"
	headers := http.Header{
		"haha": []string{"jaja"},
	}

	rb := suite.createRequest()
	rb.SetUrl(url)
	rb.SetHeaders(headers)
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	rw.Perform()
	assert.Equal(headers, mc.inReq.Header)
}

func (suite *RequestWithVerbTestSuite) TestHttpRequest_BodyPassed() {
	if !suite.verbAcceptsBody {
		return
	}

	assert := assert.New(suite.T())

	mc := &mockHttpClient{}
	url := "http://dummy-url"

	rb := suite.createRequest()
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

func (suite *RequestWithVerbTestSuite) TestHttpRequest_UnmarshallableBody() {
	if !suite.verbAcceptsBody {
		return
	}

	assert := assert.New(suite.T())

	mc := &mockHttpClient{}
	url := "http://dummy-url"

	rb := suite.createRequest()
	rb.SetUrl(url)
	rb.SetBody("kiki", unmarshallableContent{})
	rb.setHttpClient(mc)
	rw, err := rb.Build()
	assert.Nil(err)

	_, err = rw.Perform()
	assert.NotNil(err)
}

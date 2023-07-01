package connection

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHttpGetRequestTestSuite(t *testing.T) {
	suite.Run(t, &RequestWithVerbTestSuite{
		verbAcceptsBody: false,
		createRequest: func() *RequestBuilder {
			return NewHttpGetRequestBuilder()
		},
	})
}

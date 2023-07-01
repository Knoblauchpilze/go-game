package connection

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHttpDeleteRequestTestSuite(t *testing.T) {
	suite.Run(t, &RequestWithVerbTestSuite{
		verbAcceptsBody: false,
		createRequest: func() *RequestBuilder {
			return NewHttpDeleteRequestBuilder()
		},
	})
}

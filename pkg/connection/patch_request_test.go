package connection

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHttpPatchRequestTestSuite(t *testing.T) {
	suite.Run(t, &RequestWithVerbTestSuite{
		verbAcceptsBody: true,
		createRequest: func() *RequestBuilder {
			return NewHttpPatchRequestBuilder()
		},
	})
}

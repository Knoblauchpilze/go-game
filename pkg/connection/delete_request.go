package connection

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

func NewHttpDeleteRequestBuilder() *RequestBuilder {
	rb := newRequestBuilder()
	rb.setHttpRequestBuilder(buildDeleteRequest)
	return rb
}

func buildDeleteRequest(ri *requestImpl) (*http.Request, error) {
	req, err := http.NewRequest("DELETE", ri.url, nil)
	if err != nil {
		return req, errors.WrapCode(err, errors.ErrGetRequestFailed)
	}

	req.Header = ri.headers

	return req, nil
}

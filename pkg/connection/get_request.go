package connection

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

func NewHttpGetRequestBuilder() *RequestBuilder {
	rb := newRequestBuilder()
	rb.setHttpRequestBuilder(buildGetRequest)
	return rb
}

func buildGetRequest(ri *requestImpl) (*http.Request, error) {
	req, err := http.NewRequest("GET", ri.url, nil)
	if err != nil {
		return req, errors.WrapCode(err, errors.ErrGetRequestFailed)
	}

	req.Header = ri.headers

	return req, nil
}

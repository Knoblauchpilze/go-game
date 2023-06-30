package connection

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

type httpRequestBuilder func(ri *requestImpl) (*http.Request, error)

type requestImpl struct {
	url     string
	headers http.Header
	builder httpRequestBuilder
	body    interface{}
	client  HttpClient
}

func (ri *requestImpl) Perform() (*http.Response, error) {
	var resp *http.Response
	var err error

	req, err := ri.builder(ri)
	if req == nil {
		return resp, errors.New("invalid request to perform")
	}
	if err != nil {
		return resp, errors.Wrap(err, "failed to build http request")
	}

	resp, err = ri.client.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

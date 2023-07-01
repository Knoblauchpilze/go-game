package connection

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
)

func NewHttpPatchRequestBuilder() *RequestBuilder {
	rb := newRequestBuilder()
	rb.setHttpRequestBuilder(buildPatchRequest)
	return rb
}

func buildPatchRequest(ri *requestImpl) (*http.Request, error) {
	req, err := http.NewRequest("PATCH", ri.url, nil)
	if err != nil {
		return req, errors.WrapCode(err, errors.ErrGetRequestFailed)
	}

	data, err := json.Marshal(ri.body)
	if err != nil {
		return req, errors.WrapCode(err, errors.ErrPostInvalidData)
	}

	req.Header = ri.headers
	req.Body = io.NopCloser(bytes.NewReader(data))

	return req, nil
}

package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func generateRequestWithBody(body []byte) *http.Request {
	req := http.Request{}

	rdr := bytes.NewReader(body)
	req.Body = io.NopCloser(rdr)

	return &req
}

func TestGetBodyFromHttpRequestAs_NoBody(t *testing.T) {
	assert := assert.New(t)

	req := http.Request{}
	req.Body = &mockBody{}

	var in foo
	err := GetBodyFromHttpRequestAs(&req, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrFailedToGetBody))
}

func TestGetBodyFromHttpRequestAs_InvalidBody(t *testing.T) {
	assert := assert.New(t)

	var in foo

	req := generateRequestWithBody(nil)
	err := GetBodyFromHttpRequestAs(req, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))

	req = generateRequestWithBody([]byte("invalid"))
	err = GetBodyFromHttpRequestAs(req, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))
}

func TestGetBodyFromHttpRequestAs(t *testing.T) {
	assert := assert.New(t)

	in := foo{Bar: "bb", Baz: 12}
	data, _ := json.Marshal(in)
	req := generateRequestWithBody(data)

	var out foo

	err := GetBodyFromHttpRequestAs(req, &out)
	assert.Nil(err)
	assert.Equal(in.Bar, out.Bar)
	assert.Equal(in.Baz, out.Baz)
}

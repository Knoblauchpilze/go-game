package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type foo struct {
	Bar string
	Baz int
}

type mockBody struct {
	readErr  error
	closeErr error
}

func (mb mockBody) Read(p []byte) (n int, err error) {
	if mb.readErr == nil {
		return 0, io.EOF
	}

	return 0, mb.readErr
}

func (mb mockBody) Close() error {
	return mb.closeErr
}

func TestGetBodyFromHttpRequestAs_BodyReadError(t *testing.T) {
	assert := assert.New(t)

	req := http.Request{}
	req.Body = &mockBody{
		readErr: fmt.Errorf("someError"),
	}

	var in foo
	err := GetBodyFromHttpRequestAs(&req, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrFailedToGetBody))
}

func TestGetBodyFromHttpRequestAs_NilBody(t *testing.T) {
	assert := assert.New(t)

	var in foo
	req := generateRequestWithBody(nil)

	err := GetBodyFromHttpRequestAs(req, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))
}

func TestGetBodyFromHttpRequestAs_BodyNotMatchingExpectedType(t *testing.T) {
	assert := assert.New(t)

	var in foo
	req := generateRequestWithBody([]byte("invalid"))

	err := GetBodyFromHttpRequestAs(req, &in)
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

func generateRequestWithBody(body []byte) *http.Request {
	req := http.Request{}

	rdr := bytes.NewReader(body)
	req.Body = io.NopCloser(rdr)

	return &req
}

// func TestGetBodyFromHttpResponseAs_InvalidResponse(t *testing.T) {
// 	assert := assert.New(t)

// 	var in foo
// 	err := GetBodyFromHttpResponseAs(nil, &in)
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrNoResponse))

// 	resp := http.Response{
// 		StatusCode: http.StatusBadRequest,
// 	}
// 	err = GetBodyFromHttpResponseAs(&resp, &in)
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrResponseIsError))

// 	resp.StatusCode = http.StatusOK
// 	rdr := bytes.NewReader([]byte("haha"))
// 	resp.Body = io.NopCloser(rdr)

// 	err = GetBodyFromHttpResponseAs(&resp, &in)
// 	fmt.Printf("%v\n", err)
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))
// }

// func TestGetBodyFromHttpResponseAs_NoBody(t *testing.T) {
// 	assert := assert.New(t)

// 	resp := http.Response{
// 		StatusCode: http.StatusOK,
// 	}
// 	resp.Body = nil

// 	var in foo
// 	err := GetBodyFromHttpResponseAs(&resp, &in)
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrFailedToGetBody))
// }

// func TestGetBodyFromHttpResponseAs_ErrBody(t *testing.T) {
// 	assert := assert.New(t)

// 	resp := http.Response{
// 		StatusCode: http.StatusOK,
// 	}
// 	resp.Body = &mockBody{}

// 	var in foo
// 	err := GetBodyFromHttpResponseAs(&resp, &in)
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrFailedToGetBody))
// }

// func TestGetBodyFromHttpResponseAs_InvalidBody(t *testing.T) {
// 	assert := assert.New(t)

// 	var in foo

// 	resp := generateResponseWithBody(nil)
// 	err := GetBodyFromHttpResponseAs(resp, &in)
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))

// 	resp = generateResponseWithBody("invalid")
// 	err = GetBodyFromHttpResponseAs(resp, &in)
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))
// }

// func TestGetBodyFromHttpResponseAs_ErrResponse(t *testing.T) {
// 	assert := assert.New(t)

// 	in := foo{Bar: "bb", Baz: 12}
// 	resp := generateResponseWithBody(in)
// 	resp.StatusCode = http.StatusBadRequest

// 	var out foo
// 	err := GetBodyFromHttpResponseAs(resp, &out)
// 	assert.True(errors.IsErrorWithCode(err, errors.ErrResponseIsError))
// }

// func TestGetBodyFromHttpResponseAs(t *testing.T) {
// 	assert := assert.New(t)

// 	in := foo{Bar: "bb", Baz: 12}
// 	resp := generateResponseWithBody(in)

// 	var out foo

// 	err := GetBodyFromHttpResponseAs(resp, &out)
// 	assert.Nil(err)
// 	assert.Equal(in.Bar, out.Bar)
// 	assert.Equal(in.Baz, out.Baz)
// }

// func generateResponseWithBody(body interface{}) *http.Response {
// 	resp := http.Response{
// 		StatusCode: http.StatusOK,
// 	}

// 	in := NewSuccessResponse(uuid.UUID{})
// 	if body != nil {
// 		in.WithDetails(body)
// 	}

// 	data, _ := json.Marshal(in)

// 	rdr := bytes.NewReader(data)
// 	resp.Body = io.NopCloser(rdr)

// 	return &resp
// }

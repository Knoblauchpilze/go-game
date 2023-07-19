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

func TestGetBodyFromHttpResponseAs_NilResponse(t *testing.T) {
	assert := assert.New(t)

	var in foo
	err := GetBodyFromHttpResponseAs(nil, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoResponse))
}

func TestGetBodyFromHttpResponseAs_ResponseWithError(t *testing.T) {
	assert := assert.New(t)

	var in foo
	resp := http.Response{
		StatusCode: http.StatusBadRequest,
	}
	err := GetBodyFromHttpResponseAs(&resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrResponseIsError))
}

func TestGetBodyFromHttpResponseAs_ResponseOkWithNilBody(t *testing.T) {
	assert := assert.New(t)

	resp := http.Response{
		StatusCode: http.StatusOK,
		Body:       nil,
	}

	var in foo
	err := GetBodyFromHttpResponseAs(&resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrFailedToGetBody))
}

func TestGetBodyFromHttpResponseAs_ResponseNOkWithNilBody(t *testing.T) {
	assert := assert.New(t)

	resp := http.Response{
		StatusCode: http.StatusBadGateway,
		Body:       nil,
	}

	var in foo
	err := GetBodyFromHttpResponseAs(&resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrResponseIsError))
	cause := errors.Unwrap(err)
	expected := http.StatusText(http.StatusBadGateway)
	assert.Equal(expected, cause.Error())
}

func TestGetBodyFromHttpResponseAs_BodyReadError(t *testing.T) {
	assert := assert.New(t)

	resp := http.Response{
		StatusCode: http.StatusOK,
	}
	resp.Body = &mockBody{
		readErr: fmt.Errorf("someError"),
	}

	var in foo
	err := GetBodyFromHttpResponseAs(&resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrFailedToGetBody))
}

func TestGetBodyFromHttpResponseAs_ResponseWithUnexpectedBody(t *testing.T) {
	assert := assert.New(t)

	var in foo
	rdr := bytes.NewReader([]byte("haha"))
	resp := http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(rdr),
	}

	err := GetBodyFromHttpResponseAs(&resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))
}

func TestGetBodyFromHttpResponseAs_ValidBody_ResponseWithError(t *testing.T) {
	assert := assert.New(t)

	in := foo{Bar: "bb", Baz: 12}
	resp := generateResponseWithCodeAndBody(http.StatusBadRequest, in)

	var out foo
	err := GetBodyFromHttpResponseAs(resp, &out)
	assert.True(errors.IsErrorWithCode(err, errors.ErrResponseIsError))
	cause := errors.Unwrap(err)
	expected := "{\"Bar\":\"bb\",\"Baz\":12}"
	assert.Equal(expected, cause.Error())
}

func TestGetBodyFromHttpResponseAs_ValidBody_UnexpectedBody(t *testing.T) {
	assert := assert.New(t)

	resp := generateResponseWithCodeAndBody(http.StatusOK, 32)

	var out foo
	err := GetBodyFromHttpResponseAs(resp, &out)
	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))
}

func TestGetBodyFromHttpResponseAs(t *testing.T) {
	assert := assert.New(t)

	in := foo{Bar: "bb", Baz: 12}
	resp := generateResponseWithCodeAndBody(http.StatusOK, in)

	var out foo

	err := GetBodyFromHttpResponseAs(resp, &out)
	assert.Nil(err)
	assert.Equal(in.Bar, out.Bar)
	assert.Equal(in.Baz, out.Baz)
}

func generateResponseWithCodeAndBody(code int, body interface{}) *http.Response {
	mBody, _ := json.Marshal(body)

	in := ResponseTemplate{
		RequestId: dummyId,
		Status:    "someStatus",
		Details:   mBody,
	}

	data, _ := json.Marshal(in)
	rdr := bytes.NewReader(data)

	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(rdr),
	}
}

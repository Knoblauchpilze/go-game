package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var dummyId = "eb10f542-c2a8-11ed-befe-18c04d0e6a41"

func generateResponseWithBody(body interface{}) *http.Response {
	resp := http.Response{
		StatusCode: http.StatusOK,
	}

	in := NewSuccessResponse(uuid.UUID{})
	if body != nil {
		in.WithDetails(body)
	}

	data, _ := json.Marshal(in)

	rdr := bytes.NewReader(data)
	resp.Body = io.NopCloser(rdr)

	return &resp
}

func TestResponse_OK(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewSuccessResponse(id)
	resp.Write(&mrw)

	var expected = fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\"}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))
}

func TestResponse_NOK(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewErrorResponse(id)
	resp.Write(&mrw)

	var expected = fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\"}", dummyId, StatusNOK)
	assert.Equal(expected, string(mrw.data))
}

func TestResponse_Pass(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewErrorResponse(id)
	resp.Pass()
	resp.Write(&mrw)

	expected := fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\"}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))
}

func TestResponse_Fail(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewSuccessResponse(id)
	resp.Fail()
	resp.Write(&mrw)

	expected := fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\"}", dummyId, StatusNOK)
	assert.Equal(expected, string(mrw.data))
}

func TestResponse_WithDetails(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewSuccessResponse(id)
	resp.WithDetails(23)
	resp.Write(&mrw)

	expected := fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\",\"Details\":23}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))

	test := foo{
		Bar: "haha",
		Baz: -23,
	}
	resp.WithDetails(test)
	resp.Write(&mrw)

	expected = fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\",\"Details\":{\"Bar\":\"haha\",\"Baz\":-23}}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))
}

func TestResponse_WithDetails_Invalid(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewSuccessResponse(id)
	resp.WithDetails(unmarshallableContent{})
	resp.Write(&mrw)

	expected := fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\"}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))

	test := foo{
		Bar: "haha",
		Baz: -23,
	}
	resp.WithDetails(test)
	resp.Write(&mrw)

	expected = fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\",\"Details\":{\"Bar\":\"haha\",\"Baz\":-23}}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))
}

func TestResponse_WithCode(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewSuccessResponse(id)
	resp.WithCode(http.StatusOK)
	resp.Write(&mrw)

	expected := fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\"}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))
	assert.Equal(http.StatusOK, mrw.code)

	resp.WithCode(http.StatusTeapot)
	resp.Write(&mrw)

	expected = fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\"}", dummyId, StatusNOK)
	assert.Equal(expected, string(mrw.data))
	assert.Equal(http.StatusTeapot, mrw.code)
}

func TestResponse_Write(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewSuccessResponse(id)
	resp.Write(&mrw)

	expected := fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\"}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))
	assert.Equal(http.StatusOK, mrw.code)

	resp.WithDetails(12.2)
	resp.Write(&mrw)

	expected = fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\",\"Details\":12.2}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))
	assert.Equal(http.StatusOK, mrw.code)
}

func TestResponse_WriteError(t *testing.T) {
	assert := assert.New(t)

	id := uuid.MustParse(dummyId)
	mrw := mockResponseWriter{}

	resp := NewSuccessResponse(id)
	resp.WithDetails(errSomeError)
	resp.Write(&mrw)

	expected := fmt.Sprintf("{\"RequestId\":\"%s\",\"Status\":\"%v\",\"Details\":{}}", dummyId, StatusOK)
	assert.Equal(expected, string(mrw.data))
	assert.Equal(http.StatusOK, mrw.code)

}

func TestGetBodyFromHttpResponseAs_InvalidResponse(t *testing.T) {
	assert := assert.New(t)

	var in foo
	err := GetBodyFromHttpResponseAs(nil, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoResponse))

	resp := http.Response{
		StatusCode: http.StatusBadRequest,
	}
	err = GetBodyFromHttpResponseAs(&resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrResponseIsError))

	resp.StatusCode = http.StatusOK
	rdr := bytes.NewReader([]byte("haha"))
	resp.Body = io.NopCloser(rdr)

	err = GetBodyFromHttpResponseAs(&resp, &in)
	fmt.Printf("%v\n", err)
	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))
}

func TestGetBodyFromHttpResponseAs_NoBody(t *testing.T) {
	assert := assert.New(t)

	resp := http.Response{
		StatusCode: http.StatusOK,
	}
	resp.Body = nil

	var in foo
	err := GetBodyFromHttpResponseAs(&resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrFailedToGetBody))
}

func TestGetBodyFromHttpResponseAs_ErrBody(t *testing.T) {
	assert := assert.New(t)

	resp := http.Response{
		StatusCode: http.StatusOK,
	}
	resp.Body = &mockBody{}

	var in foo
	err := GetBodyFromHttpResponseAs(&resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrFailedToGetBody))
}

func TestGetBodyFromHttpResponseAs_InvalidBody(t *testing.T) {
	assert := assert.New(t)

	var in foo

	resp := generateResponseWithBody(nil)
	err := GetBodyFromHttpResponseAs(resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))

	resp = generateResponseWithBody("invalid")
	err = GetBodyFromHttpResponseAs(resp, &in)
	assert.True(errors.IsErrorWithCode(err, errors.ErrBodyParsingFailed))
}

func TestGetBodyFromHttpResponseAs_ErrResponse(t *testing.T) {
	assert := assert.New(t)

	in := foo{Bar: "bb", Baz: 12}
	resp := generateResponseWithBody(in)
	resp.StatusCode = http.StatusBadRequest

	var out foo
	err := GetBodyFromHttpResponseAs(resp, &out)
	assert.True(errors.IsErrorWithCode(err, errors.ErrResponseIsError))
}

func TestGetBodyFromHttpResponseAs(t *testing.T) {
	assert := assert.New(t)

	in := foo{Bar: "bb", Baz: 12}
	resp := generateResponseWithBody(in)

	var out foo

	err := GetBodyFromHttpResponseAs(resp, &out)
	assert.Nil(err)
	assert.Equal(in.Bar, out.Bar)
	assert.Equal(in.Baz, out.Baz)
}

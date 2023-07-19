package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestResponse_FailWithErrorAndCode(t *testing.T) {
	assert := assert.New(t)

	m := mockResponseWriter{}
	FailWithErrorAndCode(context.TODO(), nil, http.StatusBadGateway, &m)

	out, err := unmarshalExpectedResponse(m.data)
	assert.Nil(err)
	assert.Equal(statusFailure, out.Status)
	// https://www.sohamkamani.com/golang/omitempty/
	assert.Equal("", string(out.Details))
}

func TestResponse_FailWithErrorAndCode_WithDetails(t *testing.T) {
	assert := assert.New(t)

	m := mockResponseWriter{}
	inErr := errors.New("someError")
	FailWithErrorAndCode(context.TODO(), inErr, http.StatusBadGateway, &m)

	out, err := unmarshalExpectedResponse(m.data)
	assert.Nil(err)
	assert.Equal(statusFailure, out.Status)
	expected := "{\"Message\":\"someError\"}"
	assert.Equal(expected, string(out.Details))
}

func TestResponse_WriteDetails(t *testing.T) {
	assert := assert.New(t)

	m := mockResponseWriter{}
	WriteDetails(context.TODO(), nil, &m)

	out, err := unmarshalExpectedResponse(m.data)
	assert.Nil(err)
	assert.Equal(statusSuccess, out.Status)
	expected := ""
	assert.Equal(expected, string(out.Details))
}

func TestResponse_WriteDetails_WithDetails(t *testing.T) {
	assert := assert.New(t)

	m := mockResponseWriter{}
	WriteDetails(context.TODO(), 36, &m)

	out, err := unmarshalExpectedResponse(m.data)
	assert.Nil(err)
	assert.Equal(statusSuccess, out.Status)
	expected := "36"
	assert.Equal(expected, string(out.Details))
}

type mockResponseWriter struct {
	code int
	data []byte
}

func (mrw *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (mrw *mockResponseWriter) Write(out []byte) (int, error) {
	mrw.data = out
	return len(mrw.data), nil
}

func (mrw *mockResponseWriter) WriteHeader(statusCode int) {
	mrw.code = statusCode
}

func unmarshalExpectedResponse(body []byte) (ResponseTemplate, error) {
	var out ResponseTemplate
	err := json.Unmarshal(body, &out)
	return out, err
}

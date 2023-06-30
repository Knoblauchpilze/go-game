package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRequestData_FailWithErrorAndCode(t *testing.T) {
	assert := assert.New(t)

	code := http.StatusBadGateway
	inErr := errors.New("haha")
	expectedBody, err := json.Marshal(inErr)
	assert.Nil(err)

	rq := NewRequestData()
	mrw := mockResponseWriter{}
	rq.FailWithErrorAndCode(inErr, code, &mrw)

	resp, err := unmarshalExpectedResponseBody(mrw.data)
	assert.Nil(err)
	assert.Equal(code, mrw.code)
	assert.Equal("ERROR", resp.Status)
	assert.Equal(string(expectedBody), string(resp.Details))
}

func TestRequestData_WriteDetails(t *testing.T) {
	assert := assert.New(t)

	rq := NewRequestData()
	mrw := mockResponseWriter{}
	rq.WriteDetails(32, &mrw)

	resp, err := unmarshalExpectedResponseBody(mrw.data)
	assert.Nil(err)

	assert.Equal(http.StatusOK, mrw.code)
	assert.Equal("SUCCESS", resp.Status)

	var val int
	err = json.Unmarshal(resp.Details, &val)
	assert.Nil(err)
	assert.Equal(32, val)
}

func TestGetRequestDataFromContextOrFail_EmptyContext(t *testing.T) {
	assert := assert.New(t)

	mrw := mockResponseWriter{}
	req := http.Request{}

	_, res := GetRequestDataFromContextOrFail(&mrw, &req)
	assert.False(res)
	assert.Equal(http.StatusInternalServerError, mrw.code)
}

func TestGetRequestDataFromContextOrFail_ValidContext(t *testing.T) {
	assert := assert.New(t)

	mrw := mockResponseWriter{}
	req := new(http.Request)
	inRd := NewRequestData()

	ctx := context.WithValue(req.Context(), requestDataKey, inRd)
	req = req.WithContext(ctx)

	rd, res := GetRequestDataFromContextOrFail(&mrw, req)
	assert.True(res)
	assert.Equal(inRd.Id, rd.Id)
}

func TestGetRequestDataFromContextOrFail_WithContextWrapper(t *testing.T) {
	assert := assert.New(t)

	mrw := mockResponseWriter{}
	req := new(http.Request)

	next := RequestCtx(defaultHandler("haha"))
	next.ServeHTTP(&mrw, req)

	assert.Equal(http.StatusOK, mrw.code)
	resp, err := unmarshalExpectedResponseBody(mrw.data)
	assert.Nil(err)

	assert.Equal(http.StatusOK, mrw.code)
	assert.Equal("SUCCESS", resp.Status)

	var val string
	err = json.Unmarshal(resp.Details, &val)
	assert.Nil(err)
	assert.Equal("haha", val)
}

func TestGetRequestDataFromContextOrFail_WithoutContextWrapper(t *testing.T) {
	assert := assert.New(t)

	mrw := mockResponseWriter{}
	req := new(http.Request)

	next := defaultHandler("haha")
	next.ServeHTTP(&mrw, req)

	code := http.StatusInternalServerError
	assert.Equal(code, mrw.code)
	assert.Contains(string(mrw.data), http.StatusText(code))
}

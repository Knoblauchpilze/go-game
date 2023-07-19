package middlewares

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestCallerData_Host_NoReq(t *testing.T) {
	assert := assert.New(t)

	cd := callerData{}
	assert.Equal("N/A", cd.host())
}

var someHttpReq = &http.Request{
	Method:     "method",
	Host:       "someHost",
	RequestURI: "someUri",
}

func TestCallerData_Host(t *testing.T) {
	assert := assert.New(t)

	cd := callerData{
		req: someHttpReq,
	}
	assert.Equal("someHostsomeUri", cd.host())
}

var then = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

func assertThatStringStartsWith(assert *assert.Assertions, actual string, expected string) {
	assert.Greater(len(actual), len(expected))
	assert.Equal(expected, actual[:len(expected)])
}

func TestCallerData_Write_Ok(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetLoggerFuncs)

	cd := callerData{
		start: then,
		req:   someHttpReq,
	}
	m := &mockWrapResponseWriter{
		status: http.StatusOK,
	}
	var infoStr, warnStr string
	infoLog = func(ctx context.Context, format string, args ...interface{}) {
		infoStr = format
	}
	warnLog = func(ctx context.Context, format string, args ...interface{}) {
		warnStr = format
	}

	cd.write(m)
	assert.Equal("", warnStr)
	expected := "200 method from someHostsomeUri, elapsed: "
	assertThatStringStartsWith(assert, infoStr, expected)
}

func TestCallerData_Write_NOk(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetLoggerFuncs)

	cd := callerData{
		start: then,
		req:   someHttpReq,
	}
	m := &mockWrapResponseWriter{
		status: http.StatusBadRequest,
	}
	var infoStr, warnStr string
	infoLog = func(ctx context.Context, format string, args ...interface{}) {
		infoStr = format
	}
	warnLog = func(ctx context.Context, format string, args ...interface{}) {
		warnStr = format
	}

	cd.write(m)
	assert.Equal("", infoStr)
	expected := "400 method from someHostsomeUri, elapsed: "
	assertThatStringStartsWith(assert, warnStr, expected)
}

func TestCallerData_Serialize(t *testing.T) {
	assert := assert.New(t)

	cd := callerData{
		start: then,
		req:   someHttpReq,
	}
	m := &mockWrapResponseWriter{
		status: 14,
	}

	expected := "14 method from someHostsomeUri, elapsed: "
	out := cd.serialize(m)
	// Not testing for elapsed time.
	assertThatStringStartsWith(assert, out, expected)
}

func resetLoggerFuncs() {
	infoLog = logger.ScopedInfof
	warnLog = logger.ScopedWarnf
}

type mockWrapResponseWriter struct {
	status   int
	written  int
	writeErr error
	header   http.Header
}

func (m *mockWrapResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockWrapResponseWriter) Write(in []byte) (int, error) {
	return m.written, m.writeErr
}

func (m *mockWrapResponseWriter) WriteHeader(statusCode int) {}

func (m *mockWrapResponseWriter) Status() int {
	return m.status
}

func (m *mockWrapResponseWriter) BytesWritten() int {
	return m.written
}

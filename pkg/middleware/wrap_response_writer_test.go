package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapResponseWriter_Header(t *testing.T) {
	assert := assert.New(t)

	m := &mockResponseWriter{
		header: http.Header{
			"key": []string{"value1", "value2"},
		},
	}
	out := wrap(m)

	header := out.Header()
	assert.Equal(1, m.headerCalled)
	assert.Equal(m.header, header)
}

func TestWrapResponseWriter_Write(t *testing.T) {
	assert := assert.New(t)

	m := &mockResponseWriter{
		written:  26,
		writeErr: fmt.Errorf("someError"),
	}
	data := []byte{1, 2}
	out := wrap(m)

	n, err := out.Write(data)
	assert.Equal(1, m.writeCalled)
	assert.Equal(m.written, n)
	assert.Equal(m.written, out.BytesWritten())
	assert.Equal(m.writeErr, err)
	assert.Equal(data, m.data)
}

func TestWrapResponseWriter_WriteHeader(t *testing.T) {
	assert := assert.New(t)

	m := &mockResponseWriter{}
	out := wrap(m)

	out.WriteHeader(36)
	assert.Equal(1, m.writeHeaderCalled)
	assert.Equal(36, m.code)
	assert.Equal(36, out.Status())
}

type mockResponseWriter struct {
	headerCalled int
	header       http.Header

	writeCalled int
	data        []byte
	written     int
	writeErr    error

	writeHeaderCalled int
	code              int
}

func (m *mockResponseWriter) Header() http.Header {
	m.headerCalled++
	return m.header
}

func (m *mockResponseWriter) Write(out []byte) (int, error) {
	m.writeCalled++
	m.data = out
	return m.written, m.writeErr
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.writeHeaderCalled++
	m.code = statusCode
}

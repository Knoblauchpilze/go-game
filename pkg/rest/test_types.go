package rest

import (
	"fmt"
	"net/http"
)

type foo struct {
	Bar string
	Baz int
}

var errSomeError = fmt.Errorf("some error")
var errAnotherError = fmt.Errorf("another error")

type mockBody struct{}

func (mb mockBody) Read(p []byte) (n int, err error) {
	return 0, errSomeError
}

func (mb mockBody) Close() error {
	return errAnotherError
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

type unmarshallableContent struct{}

func (uc unmarshallableContent) MarshalJSON() ([]byte, error) {
	return []byte{}, errSomeError
}

package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimingCtx_WrapsResponseWriter(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDecoratorFunc)

	m := &mockHttpHandler{}
	mrw := &mockResponseWriter{}

	out := TimingCtx(m)

	out.ServeHTTP(mrw, &http.Request{})
	_, ok := m.inRespWriter.(wrapResponseWriter)
	assert.True(ok)
}

func TestTimingCtx_WritesStatus(t *testing.T) {
	assert := assert.New(t)
	// From this:
	// https://stackoverflow.com/questions/24375966/does-go-test-run-unit-tests-concurrently
	// It seems safe to have tests in the same package mutating the
	// global variables of the package as the tests are not run in
	// parallel.
	t.Cleanup(resetDecoratorFunc)

	code := http.StatusOK
	m := &mockHttpHandler{
		outCode: &code,
	}
	mrw := &mockResponseWriter{}
	var infoStr string
	infoLog = func(ctx context.Context, format string, args ...interface{}) {
		infoStr = format
	}

	out := TimingCtx(m)

	out.ServeHTTP(mrw, &http.Request{})
	expected := "200  from , elapsed: "
	assertThatStringStartsWith(assert, infoStr, expected)
}

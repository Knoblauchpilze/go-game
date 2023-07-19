package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockHttpHandler struct {
	outCode      *int
	inRespWriter http.ResponseWriter
	inReq        *http.Request
}

func (m *mockHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.inRespWriter = w
	m.inReq = r

	if m.outCode != nil {
		w.WriteHeader(*m.outCode)
	}
}

func TestRequestIdCtx_CallsDecorator(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDecoratorFunc)

	m := &mockHttpHandler{}
	decoratorCalled := false
	decoratorFunc = func(ctx context.Context, id uuid.UUID) context.Context {
		decoratorCalled = true
		return ctx
	}
	mrw := &mockResponseWriter{}

	out := RequestIdCtx(m)
	out.ServeHTTP(mrw, &http.Request{})

	assert.True(decoratorCalled)
}

func resetDecoratorFunc() {
	decoratorFunc = logger.DecorateContextWithRequestId
}

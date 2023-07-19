package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockHttpHandler struct{}

func (m *mockHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestRequestIdCtx_CallsDecorator(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDecoratorFunc)

	m := &mockHttpHandler{}
	decoratorCalled := false
	decoratorFunc = func(ctx context.Context, id uuid.UUID) context.Context {
		decoratorCalled = true
		return ctx
	}

	out := RequestIdCtx(m)

	mrw := &mockResponseWriter{}
	out.ServeHTTP(mrw, &http.Request{})

	assert.True(decoratorCalled)
}

func resetDecoratorFunc() {
	decoratorFunc = logger.DecorateContextWithRequestId
}

package middlewares

import (
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/google/uuid"
)

var decoratorFunc = logger.DecorateContextWithRequestId

func RequestIdCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		ctx := decoratorFunc(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

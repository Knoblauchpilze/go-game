package middlewares

import (
	"net/http"
	"time"
)

func TimingCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := callerData{
			start: time.Now(),
			req:   r,
		}

		rw := wrap(w)

		defer data.write(rw)

		next.ServeHTTP(rw, r)
	})
}

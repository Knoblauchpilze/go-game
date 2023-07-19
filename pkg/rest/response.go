package rest

import (
	"context"
	"net/http"
)

func FailWithErrorAndCode(ctx context.Context, err error, code int, w http.ResponseWriter) {
	b := newResponseBuilder(ctx)
	b.setCode(code)
	if err != nil {
		b.setDetails(err)
	}
	b.Write(w)
}

func WriteDetails(ctx context.Context, details interface{}, w http.ResponseWriter) {
	b := newResponseBuilder(ctx)
	if details != nil {
		b.setDetails(details)
	}
	b.Write(w)
}

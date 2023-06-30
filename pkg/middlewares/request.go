package middlewares

import (
	"context"
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/rest"
	"github.com/google/uuid"
)

type RequestData struct {
	Id       uuid.UUID
	Response rest.ResponseBuilder
}

var requestDataKey stringDataKeyType = "requestData"

func NewRequestData() RequestData {
	var out RequestData

	out.Id, _ = uuid.NewUUID()
	out.Response = rest.NewSuccessResponse(out.Id)

	return out
}

func (rd RequestData) FailWithErrorAndCode(err error, code int, w http.ResponseWriter) {
	rd.Response.WithCode(code)
	rd.Response.WithDetails(err)
	rd.Response.Write(w)
}

func (rd RequestData) WriteDetails(details interface{}, w http.ResponseWriter) {
	rd.Response.WithDetails(details)
	rd.Response.Write(w)
}

func RequestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rd := NewRequestData()

		ctx := context.WithValue(r.Context(), requestDataKey, rd)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestDataFromContextOrFail(w http.ResponseWriter, r *http.Request) (RequestData, bool) {
	ctx := r.Context()
	reqData, ok := ctx.Value(requestDataKey).(RequestData)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	return reqData, ok
}

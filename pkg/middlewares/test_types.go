package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

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

type expectedResponseBody struct {
	RequestId uuid.UUID
	Status    string
	Details   json.RawMessage
}

func unmarshalExpectedResponseBody(body []byte) (expectedResponseBody, error) {
	var out expectedResponseBody
	err := json.Unmarshal(body, &out)
	return out, err
}

func defaultHandler(msg string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rd, res := GetRequestDataFromContextOrFail(w, r)
		if res {
			rd.WriteDetails(msg, w)
		}
	})
}

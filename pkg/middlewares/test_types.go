package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

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

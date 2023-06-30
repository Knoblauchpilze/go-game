package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/auth"
	"github.com/google/uuid"
)

var errSomeError = fmt.Errorf("some error")

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

type mockAuth struct {
	isError    bool
	token      string
	expiration time.Time
}

func (ma mockAuth) GenerateToken(user uuid.UUID, password string) (auth.Token, error) {
	if ma.isError {
		return auth.Token{}, errSomeError
	}

	t := auth.Token{
		User:       user,
		Expiration: ma.expiration,
		Value:      ma.token,
	}

	return t, nil
}

func (ma mockAuth) GetToken(user uuid.UUID) (auth.Token, error) {
	return ma.GenerateToken(user, "hihi")
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

type mockServer struct {
	handler        http.Handler
	auth           mockAuth
	withRequestCtx bool
	req            *http.Request
	mrw            mockResponseWriter
}

func newMockServer(msg string) *mockServer {
	return &mockServer{
		handler: defaultHandler(msg),
		auth: mockAuth{
			isError:    false,
			token:      msg,
			expiration: time.Now(),
		},
		withRequestCtx: true,
		req: &http.Request{
			Header: make(map[string][]string),
		},
		mrw: mockResponseWriter{},
	}
}

func (ms *mockServer) withAuthorization(header string) {
	ms.req.Header["Authorization"] = []string{header}
}

func (ms *mockServer) call() {
	server := GenerateAuthenticationContext(ms.auth)(ms.handler)

	if ms.withRequestCtx {
		server = RequestCtx(server)
	}

	server.ServeHTTP(&ms.mrw, ms.req)
}

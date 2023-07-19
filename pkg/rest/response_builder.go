package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/google/uuid"
)

var statusSuccess = "SUCCESS"
var statusFailure = "ERROR"

type responseBuilder struct {
	id      uuid.UUID
	code    int
	status  string
	details json.RawMessage
}

var idUnwrapper = logger.UnwrapIdFromContext

func newResponseBuilder(ctx context.Context) *responseBuilder {
	b := responseBuilder{
		code:    http.StatusOK,
		status:  statusSuccess,
		details: nil,
	}

	id, ok := idUnwrapper(ctx)
	if !ok {
		id = uuid.New()
	}
	b.id = id

	return &b
}

func (b *responseBuilder) setCode(httpCode int) {
	b.code = httpCode
	if b.code != http.StatusOK {
		b.status = statusFailure
	} else {
		b.status = statusSuccess
	}
}

func (b *responseBuilder) setDetails(details interface{}) {
	var out []byte
	var err error

	if out, err = json.Marshal(details); err == nil {
		b.details = out
	}
}

func (b *responseBuilder) serialize() ([]byte, error) {
	toMarshal := ResponseTemplate{
		RequestId: b.id,
		Status:    b.status,
		Details:   b.details,
	}

	return json.Marshal(toMarshal)
}

func (b *responseBuilder) Write(w http.ResponseWriter) {
	out, _ := b.serialize()

	w.WriteHeader(b.code)
	w.Write(out)
}

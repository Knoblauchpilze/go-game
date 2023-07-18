package logger

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type requestIdKeyType string

const requestIdFieldName requestIdKeyType = "request"

func writeRequestIdIfFound(ctx context.Context, out io.Writer) {
	id, ok := ctx.Value(requestIdFieldName).(uuid.UUID)
	if !ok {
		return
	}

	writeColoredAndSeparateTo(id.String(), blue, out)
}

func DecorateContextWithRequestId(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, requestIdFieldName, id)
}

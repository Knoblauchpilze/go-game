package logger

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type requestIdKeyType string

const requestIdFieldName requestIdKeyType = "request"

func writeRequestIdIfFound(ctx context.Context, out io.Writer) {
	id, ok := UnwrapIdFromContext(ctx)
	if !ok {
		return
	}

	writeColoredAndSeparateTo(id.String(), blue, out)
}

func DecorateContextWithRequestId(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, requestIdFieldName, id)
}

func UnwrapIdFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(requestIdFieldName).(uuid.UUID)
	return id, ok
}

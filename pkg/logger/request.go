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

	str := writeColoredAndSeparate(id.String(), cyan)

	// Voluntarily ignoring return values.
	out.Write([]byte(str))
}

package logger

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockIoWriter struct {
	writeCalled int
	in          []byte
}

func (m *mockIoWriter) Write(p []byte) (n int, err error) {
	m.in = p
	m.writeCalled++
	return 0, nil
}

func TestWriteRequestIdIfFound_EmptyContext(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}
	ctx := context.TODO()

	writeRequestIdIfFound(ctx, m)
	assert.Equal(0, m.writeCalled)
}

func TestWriteRequestIdIfFound_ValidContext(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}
	ctx := context.TODO()
	id := uuid.New()
	ctx = context.WithValue(ctx, requestIdFieldName, id)

	writeRequestIdIfFound(ctx, m)
	assert.Equal(1, m.writeCalled)
	out := fmt.Sprintf("\033[1;36m%s\033[0m ", id)
	assert.Equal(out, string(m.in))
}

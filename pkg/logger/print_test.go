package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteAndReturn(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}
	writeAndReturnTo("", m)
	assert.Equal(1, m.writeCalled)
	assert.Equal("\n", string(m.in))

	m = &mockIoWriter{}
	writeAndReturnTo("str", m)
	assert.Equal("str\n", string(m.in))
}

func TestWriteColoredAndSeparate(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}
	writeColoredAndSeparateTo("", gray, m)
	assert.Equal(1, m.writeCalled)
	assert.Equal("\033[1;90m\033[0m ", string(m.in))

	m = &mockIoWriter{}
	writeColoredAndSeparateTo("str", yellow, m)
	assert.Equal(1, m.writeCalled)
	assert.Equal("\033[1;33mstr\033[0m ", string(m.in))
}

package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteAndReturn(t *testing.T) {
	assert := assert.New(t)

	out := writeAndReturn("")
	assert.Equal("\n", out)

	out = writeAndReturn("str")
	assert.Equal("str\n", out)
}

func TestWriteColoredAndSeparate(t *testing.T) {
	assert := assert.New(t)

	out := writeColoredAndSeparate("", gray)
	assert.Equal("\033[1;90m\033[0m ", out)

	out = writeColoredAndSeparate("str", yellow)
	assert.Equal("\033[1;33mstr\033[0m ", out)
}

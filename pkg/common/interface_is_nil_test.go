package common

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInterfaceNil_NoPanic(t *testing.T) {
	assert := assert.New(t)

	assert.NotPanics(func() {
		IsInterfaceNil(nil)
	})

	assert.NotPanics(func() {
		IsInterfaceNil(1)
	})

	var out io.ByteWriter
	assert.NotPanics(func() {
		IsInterfaceNil(out)
	})
}

func TestIsInterfaceNil(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsInterfaceNil(nil))
	assert.False(IsInterfaceNil(1))

	var out io.ByteWriter
	assert.True(IsInterfaceNil(out))
}

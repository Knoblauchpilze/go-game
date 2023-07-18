package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestWriteServiceIfFound_NotThere(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}
	f := logrus.Fields{}

	writeServiceIfFound(f, m)
	assert.Equal(0, m.writeCalled)
}

func TestWriteServiceIfFound(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}
	f := logrus.Fields{
		serviceFieldName: "service",
	}

	writeServiceIfFound(f, m)
	assert.Equal(1, m.writeCalled)
	expected := "\033[1;36m[service]\033[0m "
	assert.Equal(expected, string(m.in))
}

func TestWriteServiceIfFound_WrongType(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}
	f := logrus.Fields{
		serviceFieldName: 32,
	}

	writeServiceIfFound(f, m)
	assert.Equal(0, m.writeCalled)
}

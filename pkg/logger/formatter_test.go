package logger

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestWriteTime(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	writeTime(then, m)
	assert.Equal(1, m.writeCalled)
	expected := "\033[1;35m2009-11-17 20:34:58.651\033[0m "
	assert.Equal(expected, string(m.in))
}

func TestWriteLogLevel(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}

	writeLogLevel(logrus.DebugLevel, m)
	assert.Equal(1, m.writeCalled)
	expected := "\033[1;34m[debug]\033[0m "
	assert.Equal(expected, string(m.in))
}

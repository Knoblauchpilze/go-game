package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestColors_Trace(t *testing.T) {
	assert := assert.New(t)

	color := colorFromLogLevel(logrus.TraceLevel)
	assert.Equal(gray, color)
}

func TestColors_Debug(t *testing.T) {
	assert := assert.New(t)

	color := colorFromLogLevel(logrus.DebugLevel)
	assert.Equal(blue, color)
}

func TestColors_Info(t *testing.T) {
	assert := assert.New(t)

	color := colorFromLogLevel(logrus.InfoLevel)
	assert.Equal(green, color)
}

func TestColors_Warn(t *testing.T) {
	assert := assert.New(t)

	color := colorFromLogLevel(logrus.WarnLevel)
	assert.Equal(yellow, color)
}

func TestColors_Error(t *testing.T) {
	assert := assert.New(t)

	color := colorFromLogLevel(logrus.ErrorLevel)
	assert.Equal(red, color)
}

func TestColors_InvalidLevel(t *testing.T) {
	assert := assert.New(t)

	color := colorFromLogLevel(logrus.PanicLevel)
	assert.Equal(blue, color)
}

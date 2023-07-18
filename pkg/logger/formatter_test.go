package logger

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var then = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

func TestWriteTime(t *testing.T) {
	assert := assert.New(t)

	m := &mockIoWriter{}

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

func TestFormatter_Format(t *testing.T) {
	assert := assert.New(t)

	f := formatter{}
	e := dummyLogEntry()

	out, err := f.Format(e)
	assert.Nil(err)
	expected := "\033[1;35m2009-11-17 20:34:58.651\033[0m \033[1;32m[info]\033[0m msg\n"
	assert.Equal(expected, string(out))
}

func TestFormatter_Format_RequestId(t *testing.T) {
	assert := assert.New(t)

	f := formatter{}
	e := dummyLogEntry()
	id := uuid.New()
	e.Context = context.WithValue(e.Context, requestIdFieldName, id)

	out, err := f.Format(e)
	assert.Nil(err)
	expected := fmt.Sprintf("\033[1;35m2009-11-17 20:34:58.651\033[0m \033[1;34m%s\033[0m \033[1;32m[info]\033[0m msg\n", id)
	assert.Equal(expected, string(out))
}

func TestFormatter_Format_Service(t *testing.T) {
	assert := assert.New(t)

	f := formatter{}
	e := dummyLogEntry()
	e.Data = map[string]interface{}{
		serviceFieldName: "service",
	}

	out, err := f.Format(e)
	assert.Nil(err)
	expected := "\033[1;35m2009-11-17 20:34:58.651\033[0m \033[1;32m[info]\033[0m \033[1;36m[service]\033[0m msg\n"
	assert.Equal(expected, string(out))
}

func TestFormatter_Format_NilContext(t *testing.T) {
	assert := assert.New(t)

	f := formatter{}
	e := dummyLogEntry()
	e.Context = nil

	out, err := f.Format(e)
	assert.Nil(err)
	expected := "\033[1;35m2009-11-17 20:34:58.651\033[0m \033[1;32m[info]\033[0m msg\n"
	assert.Equal(expected, string(out))
}

func dummyLogEntry() *logrus.Entry {
	return &logrus.Entry{
		Time:    then,
		Level:   logrus.InfoLevel,
		Message: "msg",
		Context: context.TODO(),
	}
}

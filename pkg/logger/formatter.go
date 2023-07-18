package logger

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultTimeFormat = "2006-01-02 15:04:05.000"

// Inspired from:
// https://github.com/sirupsen/logrus/blob/dd1b4c2e81afc5c255f216a722b012ed26be57df/text_formatter.go
type formatter struct{}

func (t formatter) Format(logEntry *logrus.Entry) ([]byte, error) {
	out := &bytes.Buffer{}

	writeTime(logEntry.Time, out)
	writeRequestIdIfFound(logEntry.Context, out)
	writeLogLevel(logEntry.Level, out)
	writeServiceIfFound(logEntry.Data, out)
	writeAndReturnTo(logEntry.Message, out)

	return out.Bytes(), nil
}

func writeTime(t time.Time, out io.Writer) {
	timeStr := t.Format(defaultTimeFormat)
	writeColoredAndSeparateTo(timeStr, magenta, out)
}

func writeLogLevel(level logrus.Level, out io.Writer) {
	str := fmt.Sprintf("[%v]", level)
	writeColoredAndSeparateTo(str, colorFromLogLevel(level), out)
}

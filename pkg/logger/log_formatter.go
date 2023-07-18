package logger

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

const defaultTimeFormat = "2006-01-02 15:04:05.000"

const serviceFieldName = "service"

// Inspired from:
// https://github.com/sirupsen/logrus/blob/dd1b4c2e81afc5c255f216a722b012ed26be57df/text_formatter.go
type terminalFormatter struct{}

func (t terminalFormatter) Format(logEntry *logrus.Entry) ([]byte, error) {
	out := &bytes.Buffer{}

	timeStr := logEntry.Time.Format(defaultTimeFormat)
	reqStr := ""
	levelStr := fmt.Sprintf("[%v]", logEntry.Level)
	serviceStr := ""

	for field, value := range logEntry.Data {
		if field == "request" {
			if reqId, ok := value.(string); ok {
				reqStr = "[" + reqId + "]"
			}
		}
		if field == serviceFieldName {
			if service, ok := value.(string); ok {
				serviceStr = "[" + service + "]"
			}
		}
	}

	t.writeColoredAndSeparate(timeStr, magenta, out)
	if len(reqStr) > 0 {
		t.writeColoredAndSeparate(reqStr, cyan, out)
	}
	t.writeColoredAndSeparate(levelStr, colorFromLogLevel(logEntry.Level), out)
	if len(serviceStr) > 0 {
		t.writeColoredAndSeparate(serviceStr, cyan, out)
	}
	t.writeAndReturn(logEntry.Message, out)

	return out.Bytes(), nil
}

func (t terminalFormatter) writeAndReturn(msg string, out *bytes.Buffer) {
	fmt.Fprintf(out, "%s\n", msg)
}

func (t terminalFormatter) writeColoredAndSeparate(msg string, color int, out *bytes.Buffer) {
	fmt.Fprintf(out, "\033[1;%dm%s\033[0m ", color, msg)
}

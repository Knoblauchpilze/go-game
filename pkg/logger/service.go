package logger

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

const serviceFieldName = "service"

func writeServiceIfFound(fields logrus.Fields, out io.Writer) {
	serviceItf, ok := fields[serviceFieldName]
	if !ok {
		return
	}

	service, ok := serviceItf.(string)
	if !ok {
		return
	}

	str := fmt.Sprintf("[%s]", service)
	writeColoredAndSeparateTo(str, cyan, out)
}

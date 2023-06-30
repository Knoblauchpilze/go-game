package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var confLock sync.Mutex
var configuration Configuration

type Configuration struct {
	Service string
}

func Configure(config Configuration) {
	confLock.Lock()
	defer confLock.Unlock()

	logrus.SetFormatter(TerminalFormatter{})
	configuration = config
}

func withService() *logrus.Entry {
	return logrus.WithField(serviceFieldName, configuration.Service)
}

func Tracef(format string, args ...interface{}) {
	withService().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	withService().Debugf(format, args...)
}

func Printf(format string, args ...interface{}) {
	withService().Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	withService().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	withService().Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	withService().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	withService().Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	withService().Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	withService().Fatalf(format, args...)
}

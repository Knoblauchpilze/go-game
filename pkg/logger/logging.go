package logger

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

var confLock sync.Mutex
var configuration Configuration

type Configuration struct {
	Service string
	Level   logrus.Level
}

func Configure(config Configuration) {
	confLock.Lock()
	defer confLock.Unlock()

	logrus.SetFormatter(formatter{})
	configuration = config
	logrus.SetLevel(config.Level)
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

func Infof(format string, args ...interface{}) {
	withService().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	withService().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	withService().Errorf(format, args...)
}

func withServiceAndContext(ctx context.Context) *logrus.Entry {
	return withService().WithContext(ctx)
}

func ScopedTracef(ctx context.Context, format string, args ...interface{}) {
	withServiceAndContext(ctx).Tracef(format, args...)
}

func ScopedDebugf(ctx context.Context, format string, args ...interface{}) {
	withServiceAndContext(ctx).Debugf(format, args...)
}

func ScopedInfof(ctx context.Context, format string, args ...interface{}) {
	withServiceAndContext(ctx).Infof(format, args...)
}

func ScopedWarnf(ctx context.Context, format string, args ...interface{}) {
	withServiceAndContext(ctx).Warnf(format, args...)
}

func ScopedErrorf(ctx context.Context, format string, args ...interface{}) {
	withServiceAndContext(ctx).Errorf(format, args...)
}

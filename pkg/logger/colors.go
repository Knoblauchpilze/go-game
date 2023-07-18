package logger

import "github.com/sirupsen/logrus"

const (
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	magenta = 35
	cyan    = 36
	gray    = 90
)

func colorFromLogLevel(level logrus.Level) int {
	switch level {
	case logrus.TraceLevel:
		return gray
	case logrus.DebugLevel:
		return blue
	case logrus.InfoLevel:
		return green
	case logrus.WarnLevel:
		return yellow
	case logrus.ErrorLevel:
		return red
	default:
		return blue
	}
}

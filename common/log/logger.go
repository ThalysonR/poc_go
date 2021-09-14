package log

import "context"

var logger Logger

type Level int8

const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
)

func init() {
	logger = &DefaultLogger{
		level: DEBUG,
	}
}

type Logger interface {
	Debug(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
}

func GetLogger() Logger {
	loggerCopy := logger
	return loggerCopy
}

func SetLogger(newLogger Logger) {
	logger = newLogger
}

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
	Debug(...interface{})
	Error(...interface{})
	Info(...interface{})
	SetLevel(Level)
	Warn(...interface{})
	WithContext(context.Context) Logger
}

func GetLogger() Logger {
	loggerCopy := logger
	return loggerCopy
}

func SetLogger(newLogger Logger) {
	logger = newLogger
}

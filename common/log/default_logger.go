package log

import (
	"context"
	"log"
)

type DefaultLogger struct {
	level Level
}

func (d *DefaultLogger) Debug(args ...interface{}) {
	if d.level == DEBUG {
		log.Print(args...)
	}
}

func (d *DefaultLogger) Error(args ...interface{}) {
	if d.level <= ERROR {
		log.Print(args...)
	}
}

func (d *DefaultLogger) Info(args ...interface{}) {
	if d.level <= INFO {
		log.Print(args...)
	}
}

func (d *DefaultLogger) Warn(args ...interface{}) {
	if d.level <= WARN {
		log.Print(args...)
	}
}

func (d *DefaultLogger) SetLevel(level Level) {
	d.level = level
}

func (d *DefaultLogger) WithContext(context.Context) Logger {
	return d
}

package log

import (
	"context"
	"log"
)

type DefaultLogger struct {
	level Level
}

func (d *DefaultLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	if d.level == DEBUG {
		log.Print(args...)
	}
}

func (d *DefaultLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if d.level <= ERROR {
		log.Print(args...)
	}
}

func (d *DefaultLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if d.level <= INFO {
		log.Print(args...)
	}
}

func (d *DefaultLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if d.level <= WARN {
		log.Print(args...)
	}
}

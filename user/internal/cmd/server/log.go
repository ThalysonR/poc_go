package main

import (
	"context"
	"fmt"
	"time"

	"github.com/thalysonr/poc_go/common/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

type ZapLogger struct {
	ctx *context.Context
	zap *zap.Logger
}

func NewZapLogger(level log.Level) *ZapLogger {
	config := defaultConfig(level)

	logger, _ := config.Build()

	return &ZapLogger{
		zap: logger,
	}
}

func (z *ZapLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	args = append(args, z.ctx)
	z.zap.Sugar().With(spreadContext(ctx)...).Debug(mergeArgs(msg, args...)...)
}

func (z *ZapLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	args = append(args, z.ctx)
	z.zap.Sugar().With(spreadContext(ctx)...).Error(mergeArgs(msg, args...)...)
}

func (z *ZapLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	args = append(args, z.ctx)
	z.zap.Sugar().With(spreadContext(ctx)...).Info(mergeArgs(msg, args...)...)
}

func (z *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	return z
}

func (z *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	z.zap.Sugar().With(spreadContext(ctx)...).With(
		"begin", begin.Format(time.RFC3339),
		"sql", sql,
		"rows_affected", fmt.Sprintf("%d", rowsAffected),
		"error", errStr,
	).Debug("trace")
}

func (z *ZapLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	z.zap.Sugar().With(spreadContext(ctx)...).Warn(mergeArgs(msg, args...)...)
}

////////////////////////////////////////////////////////////////////////////////
///////                       AUXILIARY FUNCTIONS                        ///////
////////////////////////////////////////////////////////////////////////////////

func defaultConfig(level log.Level) zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(level)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func mergeArgs(msg string, args ...interface{}) []interface{} {
	var mergedArgs []interface{}
	mergedArgs = append(mergedArgs, msg)
	mergedArgs = append(mergedArgs, args...)
	return mergedArgs
}

func spreadContext(ctx context.Context) []interface{} {
	var args []interface{}
	// TODO: get stuff from context
	return args
}

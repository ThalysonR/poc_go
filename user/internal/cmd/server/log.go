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
	fc  func(context.Context) []interface{}
	zap *zap.Logger
}

func NewZapLogger(development bool, level log.Level, fc func(context.Context) []interface{}) *ZapLogger {
	config := defaultConfig(development, level)
	logger, _ := config.Build()

	return &ZapLogger{
		fc:  fc,
		zap: logger,
	}
}

func (z *ZapLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	z.zap.Sugar().With(z.spreadContext(ctx)...).Debug(mergeArgs(msg, args...)...)
}

func (z *ZapLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	z.zap.Sugar().With(z.spreadContext(ctx)...).Error(mergeArgs(msg, args...)...)
}

func (z *ZapLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	z.zap.Sugar().With(z.spreadContext(ctx)...).Info(mergeArgs(msg, args...)...)
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
	z.zap.Sugar().With(z.spreadContext(ctx)...).With(
		"begin", begin.Format(time.RFC3339),
		"sql", sql,
		"rows_affected", fmt.Sprintf("%d", rowsAffected),
		"error", errStr,
	).Debug("trace")
}

func (z *ZapLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	z.zap.Sugar().With(z.spreadContext(ctx)...).Warn(mergeArgs(msg, args...)...)
}

////////////////////////////////////////////////////////////////////////////////
///////                       AUXILIARY FUNCTIONS                        ///////
////////////////////////////////////////////////////////////////////////////////

func defaultConfig(development bool, level log.Level) zap.Config {
	var encoding string
	var encoder zapcore.EncoderConfig
	if development {
		encoding = "console"
		encoder = zap.NewDevelopmentEncoderConfig()
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoding = "json"
		encoder = zap.NewProductionEncoderConfig()
	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(level)),
		Development: development,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         encoding,
		EncoderConfig:    encoder,
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

func (z *ZapLogger) spreadContext(ctx context.Context) []interface{} {
	var args []interface{}
	if z.fc == nil {
		return args
	}
	if cbArgs := z.fc(ctx); len(cbArgs)%2 == 0 {
		args = cbArgs
	}
	return args
}

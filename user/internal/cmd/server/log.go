package main

import (
	"context"

	"github.com/thalysonr/poc_go/common/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func (z *ZapLogger) Debug(args ...interface{}) {
	args = append(args, z.ctx)
	z.zap.Sugar().Debug(args...)
}

func (z *ZapLogger) Error(args ...interface{}) {
	args = append(args, z.ctx)
	z.zap.Sugar().Error(args...)
}

func (z *ZapLogger) Info(args ...interface{}) {
	args = append(args, z.ctx)
	z.zap.Sugar().Info(args...)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	args = append(args, z.ctx)
	z.zap.Sugar().Warn(args...)
}

func (z *ZapLogger) SetLevel(level log.Level) {
	config := defaultConfig(level)

	logger, _ := config.Build()
	z.zap = logger
}

func (z *ZapLogger) WithContext(ctx context.Context) log.Logger {
	z.ctx = &ctx
	return z
}

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

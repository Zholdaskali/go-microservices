// internal/logger/zap_simple_adapter.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SimpleZapAdapter struct {
	sugar *zap.SugaredLogger
}

func New(level string) (Logger, error) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Устанавливаем уровень
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	baseLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &SimpleZapAdapter{
		sugar: baseLogger.Sugar(),
	}, nil
}

func (z *SimpleZapAdapter) Debug(msg string, fields ...Field) {
	z.sugar.Debugw(msg, z.fieldsToArgs(fields)...)
}

func (z *SimpleZapAdapter) Info(msg string, fields ...Field) {
	z.sugar.Infow(msg, z.fieldsToArgs(fields)...)
}

func (z *SimpleZapAdapter) Warn(msg string, fields ...Field) {
	z.sugar.Warnw(msg, z.fieldsToArgs(fields)...)
}

func (z *SimpleZapAdapter) Error(msg string, fields ...Field) {
	z.sugar.Errorw(msg, z.fieldsToArgs(fields)...)
}

func (z *SimpleZapAdapter) Fatal(msg string, fields ...Field) {
	z.sugar.Fatalw(msg, z.fieldsToArgs(fields)...)
}

func (z *SimpleZapAdapter) Debugf(format string, args ...interface{}) {
	z.sugar.Debugf(format, args...)
}

func (z *SimpleZapAdapter) Infof(format string, args ...interface{}) {
	z.sugar.Infof(format, args...)
}

func (z *SimpleZapAdapter) Errorf(format string, args ...interface{}) {
	z.sugar.Errorf(format, args...)
}

func (z *SimpleZapAdapter) With(fields ...Field) Logger {
	return &SimpleZapAdapter{
		sugar: z.sugar.With(z.fieldsToArgs(fields)...),
	}
}

func (z *SimpleZapAdapter) fieldsToArgs(fields []Field) []interface{} {
	args := make([]interface{}, 0, len(fields)*2)
	for _, field := range fields {
		args = append(args, field.Key, field.Value)
	}
	return args
}

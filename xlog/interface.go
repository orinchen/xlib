package xlog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"runtime/debug"
	"sync"
)

type ctxKey struct{}

var once sync.Once

var logger *zap.Logger

func get() *zap.Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)
		level := zap.InfoLevel
		levelEnv := os.Getenv("LOG_LEVEL")
		if levelEnv != "" {
			levelFromEnv, err := zapcore.ParseLevel(levelEnv)
			if err != nil {
				log.Println(
					fmt.Errorf("invalid level, defaulting to INFO: %w", err),
				)
			}

			level = levelFromEnv
		}

		logLevel := zap.NewAtomicLevelAt(level)

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					break
				}
			}
		}

		// log to multiple destinations (console and file)
		// extra fields are added to the JSON output alone
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, logLevel),
		)

		logger = zap.New(core)
	})
	return logger
}

func Sync() {
	_ = get().Sync()
}

func GlobalWith(fields ...zap.Field) {
	logger = get().With(fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return get().With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	get().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	get().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	get().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	get().Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	get().DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	get().Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	get().Fatal(msg, fields...)
}

// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

// WithCtx returns a copy of ctx with the Logger attached.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}

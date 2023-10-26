package xslog

import (
	"context"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"log/slog"
	"sync"
)

type ctxKey struct{}

var once sync.Once

var core zapcore.Core

var logger *slog.Logger

func InitWithZap(zapCore zapcore.Core) {
	once.Do(func() {
		core = zapCore
		logger = slog.New(zapslog.NewHandler(core, nil))
	})
}

func Sync() {
	_ = core.Sync()
}

func GlobalWith(fields ...any) {
	if logger == nil {
		return
	}
	logger = logger.With(fields...)
}

func With(fields ...any) *slog.Logger {
	if logger == nil {
		return nil
	}
	return logger.With(fields...)
}

func Debug(msg string, fields ...any) {
	if logger == nil {
		return
	}
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...any) {
	if logger == nil {
		return
	}
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...any) {
	if logger == nil {
		return
	}
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...any) {
	if logger == nil {
		return
	}
	logger.Error(msg, fields...)
}

// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return slog.Default()
}

// WithCtx returns a copy of ctx with the Logger attached.
func WithCtx(ctx context.Context, l *slog.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}

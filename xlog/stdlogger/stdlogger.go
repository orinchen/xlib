package stdlogger

import (
	"github.com/orinchen/xlib/xlog"
	"go.uber.org/zap"
)

var stdLogger *zap.Logger

func init() {
	stdLogger = xlog.CreateLogger()
}

func Sync() {
	_ = stdLogger.Sync()
}

func With(fields ...zap.Field) {
	stdLogger.With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	stdLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	stdLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	stdLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	stdLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	stdLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	stdLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	stdLogger.Fatal(msg, fields...)
}

package xslog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
)

type zapBackend struct {
	zapcore.Core
}

func (z *zapBackend) Sync() error {
	return z.Sync()
}

func creatorStdIOZapCore(conf Config) zapcore.Core {
	out := initWriter(conf.LogFile.Filename, conf.LogFile.MaxSize, conf.LogFile.MaxBackup, conf.LogFile.MaxAge)
	encoder := initEncoder(conf)
	return zapcore.NewTee(
		zapcore.NewCore(encoder, out, zap.NewAtomicLevelAt(zapLevel(conf.Level))),
	)
}

func initEncoder(conf Config) zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.TimeEncoderOfLayout(conf.TimeFormat)
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

func initWriter(filename string, maxsize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 文件位置
		MaxSize:    maxsize,   // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		Compress:   false,     // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

func InitZapBackend(conf Config) slogBackend {
	conf.Def()
	core := creatorStdIOZapCore(conf)
	sl := slog.New(zapslog.NewHandler(core, nil))
	slog.SetDefault(sl)
	return &zapBackend{core}
}

package xslog

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
)

func initSysBackend(conf Config) {
	conf.Def()
	opts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       slogLevel(conf.Level),
		ReplaceAttr: StackTraceReplaceAttr,
	}

	var fileWriter io.Writer
	if conf.LogFile != nil {
		fileWriter = &lumberjack.Logger{
			Filename:   conf.LogFile.Path,      // 日志文件的位置
			MaxSize:    conf.LogFile.MaxSize,   // 文件最大尺寸（以MB为单位）
			MaxBackups: conf.LogFile.MaxBackup, // 保留的最大旧文件数量
			MaxAge:     conf.LogFile.MaxAge,    // 保留旧文件的最大天数
			Compress:   conf.LogFile.Compress,  // 是否压缩/归档旧文件
			LocalTime:  true,                   // 使用本地时间创建时间戳
		}
	}

	switch conf.Schema {
	case "json":
		slog.SetDefault(initJsonBackend(fileWriter, opts))
	case "text":
		slog.SetDefault(initTextBackend(fileWriter, opts))
	default:
		slog.SetDefault(initTextBackend(fileWriter, opts))
	}
}

func initTextBackend(w io.Writer, opts *slog.HandlerOptions) (logger *slog.Logger) {
	if w == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(w, opts))
	}

	return logger
}

func initJsonBackend(w io.Writer, opts *slog.HandlerOptions) (logger *slog.Logger) {
	if w == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewJSONHandler(w, opts))
	}
	return logger
}

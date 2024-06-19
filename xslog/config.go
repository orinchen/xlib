package xslog

import (
	"go.uber.org/zap/zapcore"
	"log/slog"
	"strings"
)

type Config struct {
	Level      string
	TimeFormat string
	LogFile    *LogFileConfig
}

type LogFileConfig struct {
	Filename  string
	MaxSize   int
	MaxBackup int
	MaxAge    int
	Compress  bool
}

func (c *Config) Def() {
	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05"
	}
	if c.Level == "" {
		c.Level = "info"
	}
	if c.LogFile == nil {
		c.LogFile = &LogFileConfig{}
	}
	if c.LogFile.MaxSize == 0 {
		c.LogFile.MaxSize = 100
	}
	if c.LogFile.MaxBackup == 0 {
		c.LogFile.MaxBackup = 10
	}
	if c.LogFile.MaxAge == 0 {
		c.LogFile.MaxAge = 30
	}
	if c.LogFile.Compress == false {
		c.LogFile.Compress = true
	}
	if c.LogFile.Filename == "" {
		c.LogFile.Filename = "logs/app.log"
	}
}

func slogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

func zapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.ErrorLevel
	}
}

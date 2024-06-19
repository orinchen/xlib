package xslog

import (
	"go.uber.org/zap/zapcore"
	"log/slog"
	"strings"
)

type Config struct {
	Backend    string
	Level      string
	TimeFormat string         `json:",optional"`
	LogFile    *LogFileConfig `json:",optional"`
}

type LogFileConfig struct {
	Filename  string `json:",optional"`
	MaxSize   int    `json:",optional"`
	MaxBackup int    `json:",optional"`
	MaxAge    int    `json:",optional"`
	Compress  bool   `json:",optional"`
}

type slogBackend interface {
	Sync() error
}

type emptyBackend struct {
}

func (z *emptyBackend) Sync() error {
	return nil
}

func (c *Config) Def() {
	if c.Backend == "" {
		c.Backend = "color"
	}
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

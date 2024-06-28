package xslog

import (
	"log/slog"
	"strings"
)

type Config struct {
	Backend string
	Level   string
	LogFile *LogFileConfig `json:",optional"`
}

type LogFileConfig struct {
	Path      string `json:",optional"`
	MaxSize   int    `json:",optional"`
	MaxBackup int    `json:",optional"`
	MaxAge    int    `json:",optional"`
	Compress  bool   `json:",optional"`
}

func (c *Config) Def() {
	if c.Backend == "" {
		c.Backend = "dev"
	}
	if c.Level == "" {
		c.Level = "info"
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

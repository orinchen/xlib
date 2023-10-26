package xslog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
)

func CreatorStdIOZapCore() zapcore.Core {
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
	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, logLevel),
	)
}

func InitZapBackend() (core zapcore.Core) {
	core = CreatorStdIOZapCore()
	sl := slog.New(zapslog.NewHandler(core, nil))
	slog.SetDefault(sl)
	return core
}

func CreatorStdIOZapLogger() *zap.Logger {
	core := CreatorStdIOZapCore()
	return zap.New(core)
}

package xslog

import (
	"github.com/mdobak/go-xerrors"
	"log/slog"
	"testing"
)

func TestSysHandler(t *testing.T) {
	InitSlog(Config{
		Backend: "dev",
		Level:   "debug",
		LogFile: &LogFileConfig{
			Path:      "test.log",
			MaxSize:   10,
			MaxBackup: 10,
			MaxAge:    10,
			Compress:  false,
		},
	})

	slog.Error("测试错误", slog.Any("错误信息", xerrors.New("error")))
}

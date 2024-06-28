package xslog

import (
	"github.com/golang-cz/devslog"
	"log/slog"
	"os"
)

func InitDevBackend(conf Config) {
	conf.Def()
	opts := &devslog.Options{
		MaxSlicePrintSize:  4,
		SortKeys:           true,
		TimeFormat:         "[15:04:05.000]",
		NewLineAfterLog:    true,
		MaxErrorStackTrace: 6,
		HandlerOptions: &slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel(conf.Level),
		},
	}
	opts.Level = slogLevel(conf.Level)

	logger := slog.New(devslog.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return
}

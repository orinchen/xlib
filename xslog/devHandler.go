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
		TimeFormat:         conf.TimeFormat,
		NewLineAfterLog:    true,
		MaxErrorStackTrace: 4,
	}
	opts.Level = slogLevel(conf.Level)

	logger := slog.New(devslog.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

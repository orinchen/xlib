package xslog

import (
	"github.com/golang-cz/devslog"
	"log/slog"
	"os"
)

func InitDevBackend(conf Config) slogBackend {
	conf.Def()
	opts := &devslog.Options{
		MaxSlicePrintSize:  4,
		SortKeys:           true,
		TimeFormat:         conf.TimeFormat,
		NewLineAfterLog:    true,
		MaxErrorStackTrace: 4,
		HandlerOptions: &slog.HandlerOptions{
			Level: slogLevel(conf.Level),
		},
	}
	opts.Level = slogLevel(conf.Level)

	logger := slog.New(devslog.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return &emptyBackend{}
}

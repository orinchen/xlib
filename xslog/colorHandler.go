package xslog

import (
	"github.com/phsym/console-slog"
	"log/slog"
	"os"
)

func InitColorBackend(conf Config) {
	conf.Def()
	slog.SetDefault(slog.New(
		console.NewHandler(
			os.Stderr,
			&console.HandlerOptions{
				Level:      slogLevel(conf.Level),
				TimeFormat: conf.TimeFormat,
			}),
	))
}

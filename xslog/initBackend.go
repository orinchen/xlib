package xslog

import "strings"

func InitSlog(conf Config) slogBackend {
	switch strings.ToLower(conf.Backend) {
	case "zap":
		return InitZapBackend(conf)
	case "dev":
		return InitDevBackend(conf)
	case "color":
		return InitColorBackend(conf)
	default:
		return InitColorBackend(conf)
	}
}

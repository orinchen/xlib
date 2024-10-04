package xslog

import "strings"

func InitSlog(conf Config) {
	switch strings.ToLower(conf.Schema) {
	case "text":
		initSysBackend(conf)
	case "json":
		initSysBackend(conf)
	case "dev":
		InitDevBackend(conf)
	default:
		initSysBackend(conf)
	}
}

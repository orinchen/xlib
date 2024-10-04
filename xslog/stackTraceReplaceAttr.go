package xslog

import "log/slog"

func StackTraceReplaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindAny:
		switch v := a.Value.Any().(type) {
		case error:
			a.Value = fmtErr(v)
		}
	}
	return a
}
func fmtErr(err error) slog.Value {
	var groupValues []slog.Attr
	groupValues = append(groupValues, slog.String("msg", err.Error()))
	frames := getStackFrameFromPC(extractPCFromError(err, 4))
	if frames != nil {
		groupValues = append(groupValues,
			slog.Any("trace", frames),
		)
	}
	return slog.GroupValue(groupValues...)
}

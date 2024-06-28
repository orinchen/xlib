package xslog

import (
	"path/filepath"
	"reflect"
	"runtime"
)

type stackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
}

func getStackFrameFromPC(pcs []uintptr) (stackFrames []stackFrame) {
	if len(pcs) == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs[:])
	for {
		fr, more := frames.Next()
		stackFrames = append(stackFrames, stackFrame{
			Source: filepath.Join(
				filepath.Base(filepath.Dir(fr.File)),
				filepath.Base(fr.File),
			),
			Func: filepath.Base(fr.Function),
			Line: fr.Line,
		})
		if !more {
			break
		}
	}

	return
}

func extractPCFromError(err error, maxStackTrace int) (pc []uintptr) {
	if maxStackTrace == 0 {
		return nil
	}

	v := reflect.ValueOf(err)
	pc = extractPCFromPkgErrors(v, maxStackTrace)
	if len(pc) > 0 {
		return pc
	}

	pc = extractPCFromExpErrors(v, maxStackTrace)
	if len(pc) > 0 {
		return pc
	}

	return pc
}

func extractPCFromPkgErrors(v reflect.Value, maxStackTrace int) (pc []uintptr) {
	v = v.MethodByName("StackTrace")
	if !v.IsValid() {
		return nil
	}

	v = v.Call(nil)[0]
	if v.Kind() != reflect.Slice {
		return nil
	}

	// Get up to two frames from github.com/pkg/errors StackTrace.
	for i := 0; i < min(v.Len(), maxStackTrace); i++ {
		index := v.Index(i)
		if !index.CanUint() {
			return pc
		}

		pc = append(pc, uintptr(index.Uint()))
	}

	return pc
}

func extractPCFromExpErrors(v reflect.Value, maxStackTrace int) (pc []uintptr) {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	v = v.FieldByName("frame")
	if v.Kind() != reflect.Struct {
		return nil
	}

	v = v.FieldByName("frames")
	if v.Kind() != reflect.Array {
		return nil
	}

	skip := 1
	for i := skip; i < min(v.Len(), skip+maxStackTrace); i++ {
		index := v.Index(i)
		if !index.CanUint() {
			return nil
		}

		pc = append(pc, uintptr(index.Uint()))
	}

	return pc
}

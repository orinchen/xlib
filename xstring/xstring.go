package xstring

import (
	"strings"
	"unsafe"
)

func IsNilOrWhitespace(s *string) bool {
	return s == nil || IsWhitespace(*s)
}

func IsWhitespace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
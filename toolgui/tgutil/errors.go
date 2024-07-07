package tgutil

import (
	"fmt"
	"runtime"
)

func Errorf(format string, args ...interface{}) error {
	pc, _, _, ok := runtime.Caller(1)
	prefix := "unknow: "
	if ok {
		prefix = fmt.Sprintf("%s: ", runtime.FuncForPC(pc).Name())
	}
	return fmt.Errorf(prefix+format, args...)
}

func NewError(info string) error {
	pc, _, _, ok := runtime.Caller(1)
	prefix := "unknow: "
	if ok {
		prefix = fmt.Sprintf("%s: ", runtime.FuncForPC(pc).Name())
	}
	return fmt.Errorf(prefix + info)
}

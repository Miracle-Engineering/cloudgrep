package util

import (
	"errors"
	"fmt"
	"io"
	"runtime/debug"
)

var showErrorStackTrace = false

func EnableErrorStackTrace() {
	showErrorStackTrace = true
}

type stackTraceError struct {
	StackTrace string
	Err        error
}

func (e stackTraceError) Error() string {
	return e.Err.Error()
}
func (e stackTraceError) Unwrap() error { return e.Err }

func GetStackTrace(e error) string {
	var cgerr stackTraceError
	if errors.As(e, &cgerr) {
		return cgerr.StackTrace
	}
	return ""
}

func AddStackTrace(e error) error {
	if showErrorStackTrace {
		return stackTraceError{
			StackTrace: string(debug.Stack()),
			Err:        e,
		}
	}
	return e
}

func PrintStacktrace(e error, w io.Writer) {
	var err error
	stackTrace := GetStackTrace(e)
	if stackTrace != "" {
		_, err = fmt.Fprintf(w, "Stack Trace:\n%v\n%v", stackTrace, e)
	}
	if err != nil {
		panic(err)
	}
}

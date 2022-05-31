package testingutil

import (
	"reflect"
	"strings"
)

// FuncSignature returns the signature of a function as a string, useful for error messages.
func FuncSignature(value any) string {
	var t reflect.Type
	t, ok := value.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(value)
	}

	// Source: https://stackoverflow.com/a/54129236
	if t.Kind() != reflect.Func {
		return "<not a function>"
	}

	buf := strings.Builder{}
	buf.WriteString("func (")
	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(t.In(i).String())
	}
	buf.WriteString(")")
	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			buf.WriteString(" (")
		} else {
			buf.WriteString(" ")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(t.Out(i).String())
		}
		if numOut > 1 {
			buf.WriteString(")")
		}
	}

	return buf.String()
}

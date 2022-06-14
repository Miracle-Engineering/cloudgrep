package testingutil

import (
	"fmt"
	"path"
	"reflect"
)

// TypeStr is a convenience function to get the fully qualified type identifier for a value
func TypeStr(v any) string {
	t := reflect.TypeOf(v)
	return fmt.Sprintf("%v/%v", path.Dir(t.PkgPath()), t.String())
}

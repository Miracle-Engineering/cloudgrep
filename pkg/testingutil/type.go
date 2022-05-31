package testingutil

import (
	"fmt"
	"path"
	"reflect"
)

func TypeStr(v any) string {
	t := reflect.TypeOf(v)
	return fmt.Sprintf("%v/%v", path.Dir(t.PkgPath()), t.String())
}

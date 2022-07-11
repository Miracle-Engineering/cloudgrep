package testingutil

import (
	"fmt"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TypeStr is a convenience function to get the fully qualified type identifier for a value
func TypeStr(v any) string {
	t := reflect.TypeOf(v)
	return fmt.Sprintf("%v/%v", path.Dir(t.PkgPath()), t.String())
}

// TestingTB is an interface wrapper around testing.TB reduced to what the funcs in testingutil need (to assist with tesing this package).
// Satisfies the assert.TestingT and require.TestingT interfaces
type TestingTB interface {
	Errorf(string, ...any)
	Fatalf(string, ...any)
	Helper()
	FailNow()
}

var _ TestingTB = &testing.T{}
var _ TestingTB = &testing.B{}
var _ assert.TestingT = TestingTB(&testing.T{})
var _ require.TestingT = TestingTB(&testing.T{})

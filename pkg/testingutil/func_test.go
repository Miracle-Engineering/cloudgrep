package testingutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncSignature(t *testing.T) {
	tests := []struct {
		name     string
		f        any
		expected string
	}{
		{
			f:        func() {},
			expected: "func ()",
		},
		{
			f:        1,
			expected: "<not a function>",
		},
		{
			name:     "reflect.Type",
			f:        reflect.TypeOf(func() {}),
			expected: "func ()",
		},
		{
			name:     "multiple in/out",
			f:        func(a string, b int) (c bool, d error) { return },
			expected: "func (string, int) (bool, error)",
		},
		{
			name:     "single in/out",
			f:        func(b int) error { return nil },
			expected: "func (int) error",
		},
	}

	for _, test := range tests {
		name := test.name
		if name == "" {
			name = test.expected
		}

		t.Run(name, func(t *testing.T) {
			actual := FuncSignature(test.f)
			assert.Equal(t, test.expected, actual)
		})
	}
}

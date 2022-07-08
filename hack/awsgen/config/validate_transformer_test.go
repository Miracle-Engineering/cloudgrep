package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformer_Validate(t *testing.T) {
	type test struct {
		expr  string
		valid bool
	}
	good := func(e string) test {
		return test{expr: e, valid: true}
	}
	bad := func(e string) test {
		return test{expr: e, valid: false}
	}
	tests := []test{
		good("foo"),
		bad("1, 2"),
		good("generic[%type]"),
		bad("generic[%type]."),
	}

	for _, testCase := range tests {
		t.Run(testCase.expr, func(t *testing.T) {
			transformer := Transformer{Expr: testCase.expr}
			errs := transformer.Validate()
			errStrs := errorsToStrings(errs)

			var expected []string

			if !testCase.valid {
				expected = []string{
					"not a valid Go expression",
				}
			}
			assert.ElementsMatch(t, expected, errStrs)
		})
	}
}

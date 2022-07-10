package config

import (
	"errors"
	"strings"
)

func (t Transformer) Validate() []error {
	var errs []error

	errs = append(errs, validateFuncs(t,
		validateTransformerExpression,
	)...)

	return errs
}

func validateTransformerExpression(t Transformer) []error {
	expr := t.Expr
	expr = strings.TrimSpace(expr)

	if expr == "" {
		return []error{
			errors.New("expr is required"),
		}
	}

	if t.IsGeneric() {
		expr = strings.ReplaceAll(expr, TransformerTypePlaceholder, "types.T")
	}

	if !isValidGoExpression(expr) {
		return []error{
			errors.New("not a valid Go expression"),
		}
	}

	return nil
}

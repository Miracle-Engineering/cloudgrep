package config

import (
	"go/token"
	"strings"
)

func isValidExportedIdentifier(id string) bool {
	return token.IsIdentifier(id) && token.IsExported(id)
}

func isValidTypeRef(typ string) bool {
	if typ == "string" {
		return true
	}

	parts := strings.SplitN(typ, ".", 2)
	if len(parts) == 1 {
		// We expect type refs to point at types defined in other packages
		// (specifically in some package in the aws SDK)
		return false
	}

	for idx, part := range parts {
		if !token.IsIdentifier(part) {
			return false
		}

		if idx > 0 && !token.IsExported(part) {
			// The first part is a package
			// the second part must be an exported identifier
			return false
		}
	}

	return true
}

func validateFuncs[T any](val T, funcs ...func(T) []error) []error {
	var errs []error
	for _, f := range funcs {
		funcErrs := f(val)
		errs = append(errs, funcErrs...)
	}

	return errs
}

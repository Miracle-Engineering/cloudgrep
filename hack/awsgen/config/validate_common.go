package config

import "fmt"

func validateExportedIdentifier(name, val string) []error {
	if !isValidExportedIdentifier(val) {
		return []error{fmt.Errorf("%s not a valid Go exported identifier: %s", name, val)}
	}

	return nil
}

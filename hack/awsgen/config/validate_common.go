package config

import "fmt"

func validateExportedIdentifier(name, val string) []error {
	if !isValidExportedIdentifier(val) {
		return []error{fmt.Errorf("%s not a valid Go exported identifier: %s", name, val)}
	}

	return nil
}

func validateInputOverrides(io InputOverrides) []error {
	var errs []error

	for field, funcName := range io.FieldFuncs {
		errs = append(errs, validateExportedIdentifier("fieldFuncs[]", field)...)

		if !isValidNameRef(funcName) {
			errs = append(errs, fmt.Errorf("fieldFuncs[%s] is not a valid func ref: %s", field, funcName))
		}
	}

	for idx, funcName := range io.FullFuncs {
		if !isValidNameRef(funcName) {
			errs = append(errs, fmt.Errorf("fullFuncs[%d] is not a valid func ref: %s", idx, funcName))
		}
	}

	setErrContextExtraPrepend("inputOverrides", errs)

	return errs
}

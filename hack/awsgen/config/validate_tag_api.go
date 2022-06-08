package config

import (
	"errors"
	"fmt"
)

func (api GetTagAPI) Validate() []error {
	var errs []error

	if api.Has() {
		errs = append(errs, validateFuncs(api,
			validateTagAPICall,
			validateTagAPIResourceType,
			validateTagAPIInputIDField,
			validateTagAPIOutputKey,
			validateTagAPITagField,
			validateTagAPIAllowedAPIErrorCodes,
			validateTagAPIUnset,
		)...)
	} else {
		errs = append(errs, validateFuncs(api,
			validateTagAPIUnset,
		)...)
	}

	return errs
}

func validateTagAPICall(api GetTagAPI) []error {
	return validateAPICall(api.Call)
}

func validateTagAPIResourceType(api GetTagAPI) []error {
	if api.ResourceType == "" {
		return []error{errors.New("resourceType required")}
	}

	return validateExportedIdentifier("resourceType", api.ResourceType)
}

func validateTagAPIInputIDField(api GetTagAPI) []error {
	errs := api.InputIDField.Validate()
	setErrContextExtraPrepend("inputIDField", errs)
	return errs
}

func validateTagAPIOutputKey(api GetTagAPI) []error {
	return api.OutputKey.Validate("outputKey")
}

func validateTagAPITagField(api GetTagAPI) []error {
	if api.TagField == nil {
		return nil
	}

	errs := api.TagField.Validate()
	setErrContextExtraPrepend("tags", errs)

	return errs
}

func validateTagAPIAllowedAPIErrorCodes(api GetTagAPI) []error {
	var errs []error
	for idx, code := range api.AllowedAPIErrorCodes {
		// Error codes are just strings, but they should have the same
		// formatting rules as exported Go identifiers
		if !isValidExportedIdentifier(code) {
			errs = append(errs, fmt.Errorf("allowedApiErrorCodes[%d]: not a valid error code: %s", idx, code))
		}
	}

	return errs
}

func validateTagAPIUnset(api GetTagAPI) []error {
	if api.Has() {
		return nil
	}

	var errs []error

	msgFmt := "expected `call` to be set when %s is set"

	add := func(name string) {
		errs = append(errs, fmt.Errorf(msgFmt, name))
	}

	if api.ResourceType != "" {
		add("type")
	}

	if !api.InputIDField.Zero() {
		add("inputIDField")
	}

	if !api.OutputKey.Empty() {
		add("outputKey")
	}

	// api.TagField != nil already validated by type

	if len(api.AllowedAPIErrorCodes) > 0 {
		add("allowedApiErrorCodes")
	}

	return errs
}

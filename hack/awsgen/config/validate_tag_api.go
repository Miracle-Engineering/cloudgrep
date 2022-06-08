package config

import (
	"errors"
	"fmt"
)

func (api GetTagsAPI) Validate() []error {
	var errs []error

	if api.Has() {
		errs = append(errs, validateFuncs(api,
			validateTagAPICall,
			validateTagAPIResourceType,
			validateTagAPIInputIDField,
			validateTagAPITags,
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

func validateTagAPICall(api GetTagsAPI) []error {
	return validateAPICall(api.Call)
}

func validateTagAPIResourceType(api GetTagsAPI) []error {
	if api.ResourceType == "" {
		return []error{errors.New("resourceType required")}
	}

	return validateExportedIdentifier("resourceType", api.ResourceType)
}

func validateTagAPIInputIDField(api GetTagsAPI) []error {
	errs := api.InputIDField.Validate()
	setErrContextExtraPrepend("inputIDField", errs)
	return errs
}

func validateTagAPITags(api GetTagsAPI) []error {
	if api.Tags == nil {
		return nil
	}

	errs := api.Tags.Validate()
	setErrContextExtraPrepend("tags", errs)

	return errs
}

func validateTagAPIAllowedAPIErrorCodes(api GetTagsAPI) []error {
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

func validateTagAPIUnset(api GetTagsAPI) []error {
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

	// api.TagField != nil already validated by type

	if len(api.AllowedAPIErrorCodes) > 0 {
		add("allowedApiErrorCodes")
	}

	return errs
}

package config

import (
	"errors"
	"fmt"
)

func (api ListAPI) Validate() []error {
	var errs []error

	errs = append(errs, validateFuncs(api,
		validateListAPICall,
		validateListAPIOutputKey,
		validateListAPIIDField,
		validateListAPITagField,
		validateListAPIInputOverrides,
	)...)

	return errs
}

func validateListAPICall(api ListAPI) []error {
	return validateAPICall(api.Call)
}

func validateListAPIOutputKey(api ListAPI) []error {
	var errs []error

	if len(api.OutputKey) == 0 {
		errs = append(errs, errors.New("outputKey empty"))
	}

	for idx, key := range api.OutputKey {
		ref := fmt.Sprintf("outputKey[%d]", idx)
		if len(key) == 0 {
			errs = append(errs, fmt.Errorf("%s is an empty string", ref))
		} else {
			errs = append(errs, validateExportedIdentifier(ref, key)...)
		}
	}

	return errs
}

func validateListAPIIDField(api ListAPI) []error {
	var errs []error

	f := api.IDField
	if f.SliceType != "" {
		errs = append(errs, errors.New("sliceType cannot be set"))

		// Make sure Field.Validate doesn't do checks on the SliceType field
		f.SliceType = ""
	}

	errs = append(errs, f.Validate()...)

	setErrContextExtraPrepend("id", errs)

	return errs
}

func validateListAPITagField(api ListAPI) []error {
	if api.Tags == nil {
		return nil
	}

	var errs []error

	if api.Tags.Style == "map" {
		// TODO: We should support this. Wait until we have such a resource?
		errs = append(errs, errors.New("map style tags not yet supported on listApi"))
	} else if api.Tags.Style == "" {
		// While we only support the "struct" style in list,
		// Default to struct if not already set
		api.Tags.Style = "struct"
	}

	errs = append(errs, api.Tags.Validate()...)

	setErrContextExtraPrepend("tags", errs)

	return errs
}

func validateListAPIInputOverrides(api ListAPI) []error {
	return validateInputOverrides(api.InputOverrides)
}

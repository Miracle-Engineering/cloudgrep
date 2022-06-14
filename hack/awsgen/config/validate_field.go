package config

import (
	"errors"
	"fmt"
)

type fieldValidationOpts struct {
	slicesProhibited  bool
	pointerProhibited bool
}

func (f Field) Validate() []error {
	return f.validate(fieldValidationOpts{})
}

func (f Field) validate(opts fieldValidationOpts) []error {
	var errs []error

	if opts.pointerProhibited && f.Pointer {
		errs = append(errs, errors.New("pointer not supported"))
		f.Pointer = false
	}

	if opts.slicesProhibited && f.SliceType != "" {
		errs = append(errs, errors.New("sliceType not supported"))
		f.SliceType = ""
	}

	if f.SliceType != "" && f.Pointer {
		errs = append(errs, errors.New("pointer and sliceType are mutually exclusive"))
	}

	if f.SliceType != "" && !isValidTypeRef(f.SliceType) {
		errs = append(errs, fmt.Errorf("sliceType does not refer to a valid type: %s", f.SliceType))
	}

	if f.Name == "" {
		errs = append(errs, errors.New("name required"))
	} else {
		errs = append(errs, validateExportedIdentifier("name", f.Name)...)
	}

	return errs
}

func (nf NestedField) Validate(ctx string) []error {
	return nf.validate(ctx, fieldValidationOpts{})
}

// ValidateSimple is like Validate, but additionally enforces that sliceType and pointer are not set
func (nf NestedField) ValidateSimple(ctx string) []error {
	opts := fieldValidationOpts{
		pointerProhibited: true,
		slicesProhibited:  true,
	}

	return nf.validate(ctx, opts)
}

func (nf NestedField) validate(ctx string, opts fieldValidationOpts) []error {
	var errs []error

	if nf.Empty() {
		errs = append(errs, fmt.Errorf("%s cannot be empty", ctx))
	}

	for idx, f := range nf {
		fieldErrs := f.validate(opts)
		ctx := fmt.Sprintf("%s[%d]", ctx, idx)
		setErrContextExtraPrepend(ctx, fieldErrs)
		errs = append(errs, fieldErrs...)
	}

	return errs
}

func (f TagField) Validate() []error {
	var errs []error

	errs = append(errs, f.Field.Validate("field")...)

	var v func(string, string)
	if f.Style == "" {
		errs = append(errs, errors.New("style required"))
		return errs
	} else if f.Style == "map" {
		v = func(name, val string) {
			if val != "" {
				errs = append(errs, fmt.Errorf("%s must not be set with style=map", name))
			}
		}

		if f.Pointer {
			errs = append(errs, errors.New("pointer not supported with style=map"))
		}
	} else if f.Style == "struct" {
		v = func(name, val string) {
			if val == "" {
				errs = append(errs, fmt.Errorf("%s required with style=struct", name))
			} else {
				errs = append(errs, validateExportedIdentifier(name, val)...)
			}
		}
	} else {
		errs = append(errs, fmt.Errorf("unknown style: %s", f.Style))
		return errs
	}

	v("key", f.Key)
	v("value", f.Value)

	// setErrContextExtraPrepend("tags", errs)
	return errs
}

package config

import (
	"errors"
	"fmt"
	"strings"
)

type validationErr struct {
	svc          *Service
	typ          *Type
	extraContext string
	wrapped      error
}

func (err validationErr) Error() string {
	buf := strings.Builder{}

	if err.svc != nil {
		buf.WriteString(fmt.Sprintf("service '%s': ", err.svc.Name))
	}

	if err.typ != nil {
		buf.WriteString(fmt.Sprintf("type '%s': ", err.typ.Name))
	}

	if err.extraContext != "" {
		buf.WriteString(fmt.Sprintf("%s: ", err.extraContext))
	}

	if err.wrapped != nil {
		buf.WriteString(err.wrapped.Error())
	} else {
		// fallback error message
		buf.WriteString("validation error")
	}

	return buf.String()
}

func (err validationErr) Unwrap() error {
	return err.wrapped
}

func typeValidationError(typ Type, wrapped error) error {
	return validationErr{typ: &typ, wrapped: wrapped}
}

func typeValidationErrorS(typ Type, msg string) error {
	return typeValidationError(typ, errors.New(msg))
}

func setErrContextService(svc Service, errs []error) {
	for idx, err := range errs {
		var validErr validationErr
		if errors.As(err, &validErr) {
			validErr.svc = &svc
		} else {
			validErr = validationErr{svc: &svc, wrapped: err}
		}

		errs[idx] = validErr
	}
}

func setErrContextType(typ Type, errs []error) {
	for idx, err := range errs {
		var validErr validationErr
		if errors.As(err, &validErr) {
			validErr.typ = &typ
		} else {
			validErr = validationErr{typ: &typ, wrapped: err}
		}

		errs[idx] = validErr
	}
}

func setErrContextExtraPrepend(extraContext string, errs []error) {
	for idx, err := range errs {
		var validErr validationErr
		if errors.As(err, &validErr) {
			if validErr.extraContext == "" {
				validErr.extraContext = extraContext
			} else {
				validErr.extraContext = fmt.Sprintf("%s: %s", extraContext, validErr.extraContext)
			}
		} else {
			validErr = validationErr{extraContext: extraContext, wrapped: err}
		}

		errs[idx] = validErr
	}
}

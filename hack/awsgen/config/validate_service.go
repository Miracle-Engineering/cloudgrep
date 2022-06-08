package config

import (
	"errors"
	"fmt"
	"regexp"
)

func (s Service) Validate() []error {
	var errs []error

	errs = append(errs, validateFuncs(
		s,
		validateServiceName,
		validateServicePackageName,
		validateServiceTypesUnique,
	)...)

	for _, typ := range s.Types {
		typErrs := typ.Validate()
		errs = append(errs, typErrs...)
	}

	setErrContextService(s, errs)

	return errs
}

const serviceNameRegex = "^[a-z][a-z0-9]*$"

func validateServiceName(service Service) []error {
	if match, _ := regexp.MatchString(serviceNameRegex, service.Name); !match {
		return []error{
			errors.New("name not valid"),
		}
	}

	return nil
}

func validateServicePackageName(service Service) []error {
	if match, _ := regexp.MatchString(serviceNameRegex, service.ServicePackage); !match {
		return []error{
			fmt.Errorf("servicePackage not valid: %s", service.ServicePackage),
		}
	}

	return nil
}

func validateServiceTypesUnique(service Service) []error {
	var errs []error

	// Use value to track if we have already emitted an error for that type name
	names := make(map[string]bool)

	for _, typ := range service.Types {
		name := typ.Name
		errored, has := names[name]

		if has && !errored {
			names[name] = true
			errs = append(errs, typeValidationErrorS(typ, "duplicate type name"))
		} else if !has {
			names[name] = false
		}
	}

	return errs
}

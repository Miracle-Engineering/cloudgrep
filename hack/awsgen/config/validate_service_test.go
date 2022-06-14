package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Validate(t *testing.T) {
	svc := Service{
		Name:           "Foo",
		ServicePackage: "Bar",
		Types: []Type{
			{Name: "Spam"},
			{Name: "Spam"},
			{Name: "Spam"},
		},
	}

	expected := []string{
		"service 'Foo': name not valid",
		"service 'Foo': servicePackage not valid: Bar",
		"service 'Foo': type 'Spam': duplicate type name",
	}

	svcErrs := svc.Validate()
	errStrs := serviceValidateRemoveTypeErrors(svc, svcErrs)

	assert.ElementsMatch(t, expected, errStrs)
}

func serviceValidateRemoveTypeErrors(svc Service, errs []error) []string {
	subErrs := svc.subValidate()
	setErrContextService(svc, subErrs)
	subErrStrs := errorsToStrings(subErrs)

	svcErrStrings := errorsToStrings(errs)

	return sliceDiff(svcErrStrings, subErrStrs)
}

func errorsToStrings(errs []error) []string {
	out := make([]string, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}

		out = append(out, err.Error())
	}

	return out
}

func sliceDiff[T comparable](a, b []T) []T {
	var out []T
	skipItems := make(map[T]struct{})

	for _, item := range b {
		skipItems[item] = struct{}{}
	}

	for _, item := range a {
		if _, has := skipItems[item]; has {
			continue
		}

		out = append(out, item)
	}

	return out
}

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTagsAPI_Validate_unset(t *testing.T) {
	api := GetTagsAPI{
		ResourceType:         "foo",
		InputIDField:         Field{Name: "bar"},
		AllowedAPIErrorCodes: []string{"spam"},
	}
	expected := []string{
		"expected `call` to be set when type is set",
		"expected `call` to be set when inputIDField is set",
		"expected `call` to be set when allowedApiErrorCodes is set",
	}

	assertGetTagsApiErrors(t, api, expected)
}

func assertGetTagsApiErrors(t *testing.T, api GetTagsAPI, expected []string) bool {
	t.Helper()

	errs := api.Validate()
	errStrs := errorsToStrings(errs)

	return assert.ElementsMatch(t, expected, errStrs)
}

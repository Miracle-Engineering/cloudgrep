package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAPI_Validate_1(t *testing.T) {
	api := ListAPI{
		Call: "foo",
		Tags: &TagField{
			Style: "struct",
			Value: "bar",
		},
	}
	expected := []string{
		"call is not a valid Go exported identifier: foo",
		"outputKey cannot be empty",
		"id: name required",
		"tags: field cannot be empty",
		"tags: key required with style=struct",
		"tags: value is not a valid Go exported identifier: bar",
	}

	assertListApiErrors(t, api, expected)
}

func TestListAPI_Validate_2(t *testing.T) {
	api := ListAPI{
		Call: "Foo",
		OutputKey: NestedField{
			Field{Name: "Foo"},
			Field{Name: "bar", SliceType: "string", Pointer: true},
			Field{Name: ""},
		},
		IDField: Field{
			Name:      "spam",
			SliceType: "string",
		},
		Tags: &TagField{
			Style: "map",
		},
	}

	expected := []string{
		"outputKey[1]: name is not a valid Go exported identifier: bar",
		"outputKey[1]: pointer not supported",
		"outputKey[1]: sliceType not supported",
		"outputKey[2]: name required",
		"id: sliceType cannot be present",
		"id: name is not a valid Go exported identifier: spam",
		"tags: map style tags not yet supported on listApi",
	}

	assertListApiErrors(t, api, expected)
}

func assertListApiErrors(t *testing.T, api ListAPI, expected []string) bool {
	t.Helper()

	errs := api.Validate()
	errStrs := errorsToStrings(errs)

	return assert.ElementsMatch(t, expected, errStrs)
}

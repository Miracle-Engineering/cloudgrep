package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_Validate_nameInvalid(t *testing.T) {
	typ := Type{
		Name:       "foo",
		GetTagsAPI: GetTagsAPI{Tags: &TagField{}},
	}

	expected := []string{
		"type 'foo': name not valid",
		"type 'foo': listApi.tags must be set when not configuring getTagsApi.call",
		"type 'foo': getTagsApi.tags must not be set when not configuring getTagsApi.call",
	}

	errs := typ.Validate()
	errStrs := typeValidateRemoveApiErrors(typ, errs)

	assert.ElementsMatch(t, expected, errStrs)
}

func TestType_Validate_missingGetTagsApiTags(t *testing.T) {
	typ := Type{
		Name:       "Foo",
		ListAPI:    ListAPI{Tags: &TagField{}},
		GetTagsAPI: GetTagsAPI{Call: "Foo"},
	}

	expected := []string{
		"type 'Foo': listApi.tags must not be set when configuring getTagsApi.call",
		"type 'Foo': getTagsApi.tags must be set when configuring getTagsApi.call",
	}

	errs := typ.Validate()
	errStrs := typeValidateRemoveApiErrors(typ, errs)

	assert.ElementsMatch(t, expected, errStrs)
}

func TestType_Validate_transformersInvalid(t *testing.T) {
	typ := Type{
		Name:    "Foo",
		ListAPI: ListAPI{Tags: &TagField{}},
		Transformers: []Transformer{
			{
				Name: "bar",
				Expr: ".",
			},
			{
				Name: "2",
				Expr: "",
			},
			{
				Name: "spam",
				Expr: "%type",
			},
			{
				Name: "bar",
				Expr: "",
			},
			{
				Name: "",
				Expr: "!",
			},
		},
	}

	expected := []string{
		"type 'Foo': transformers[bar]: not a valid Go expression",
		"type 'Foo': transformers[1]: name cannot be an int",
		"type 'Foo': transformers[3]: name is duplicated: bar",
		"type 'Foo': transformers[3]: expr is required",
		"type 'Foo': transformers[4]: not a valid Go expression",
	}

	errs := typ.Validate()
	errStrs := typeValidateRemoveApiErrors(typ, errs)

	assert.ElementsMatch(t, expected, errStrs)
}

func typeValidateRemoveApiErrors(typ Type, errs []error) []string {
	subErrs := typ.subValidate()
	setErrContextType(typ, subErrs)
	subErrStrs := errorsToStrings(subErrs)

	typeErrStrs := errorsToStrings(errs)

	return sliceDiff(typeErrStrs, subErrStrs)
}

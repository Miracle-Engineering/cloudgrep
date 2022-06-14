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

func typeValidateRemoveApiErrors(typ Type, errs []error) []string {
	subErrs := typ.subValidate()
	setErrContextType(typ, subErrs)
	subErrStrs := errorsToStrings(subErrs)

	typeErrStrs := errorsToStrings(errs)

	return sliceDiff(typeErrStrs, subErrStrs)
}

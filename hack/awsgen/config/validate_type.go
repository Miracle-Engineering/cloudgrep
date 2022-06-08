package config

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func (t Type) Validate() []error {
	var errs []error

	errs = append(errs, validateFuncs(t,
		validateTypeName,
		validateTypeDescription,
		validateTypeTags,
	)...)

	listErrs := t.ListAPI.Validate()
	setErrContextExtraPrepend("listApi", listErrs)
	errs = append(errs, listErrs...)

	tagErrs := t.GetTagsAPI.Validate()
	setErrContextExtraPrepend("getTagsApi", tagErrs)
	errs = append(errs, tagErrs...)

	setErrContextType(t, errs)

	return errs
}

const typeNameRegex = "^[A-Z][a-zA-Z0-9]*$"

func validateTypeName(typ Type) []error {
	if match, _ := regexp.MatchString(typeNameRegex, typ.Name); !match {
		return []error{
			typeValidationErrorS(typ, "name not valid"),
		}
	}

	return nil
}

func validateTypeDescription(typ Type) []error {
	if strings.TrimSpace(typ.Description) == "" {
		return []error{
			typeValidationErrorS(typ, "type description cannot be empty"),
		}
	}

	return nil
}

func validateTypeTags(typ Type) []error {
	var errs []error

	const name = "getTagsApi"

	m := map[string]*TagField{
		"listApi": typ.ListAPI.Tags,
		name:      typ.GetTagsAPI.TagField,
	}

	hasTags := typ.GetTagsAPI.Has()

	keys := maps.Keys(m)
	slices.Sort(keys)
	for _, key := range keys {
		hasRef := m[key] != nil
		tagRef := key == name

		msgFmt := "%s.tags %s be set %s configuring %s.call"

		configured := "when"
		if !hasTags {
			configured = "when not"
		}

		expectation := "must"
		expected := true
		if hasTags != tagRef {
			expectation = "must not"
			expected = false
		}

		if hasRef != expected {
			msg := fmt.Sprintf(msgFmt, key, expectation, configured, name)
			errs = append(errs,
				typeValidationErrorS(typ, msg),
			)
		}
	}

	return errs
}

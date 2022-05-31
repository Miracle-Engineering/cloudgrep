package testingutil

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
)

const TestRegion = "us-east-1"
const TestTag = "test"

// AssertResourceCount asserts that there is a specific number of given resources with the "test" tag.
// If tagValue is not an empty string, it also filters on resources that have the "test" tag with that value.
func AssertResourceCount(t testing.TB, resources []model.Resource, tagValue string, count int) {
	t.Helper()
	if tagValue == "" {
		resources = ResourceFilterTagKey(resources, TestTag)
	} else {
		resources = ResourceFilterTagKeyValue(resources, TestTag, tagValue)
	}

	assert.Lenf(t, resources, count, "expected %d resource(s) with tag %s=%s", count, TestTag, tagValue)
}

// ResourceFilterTagKey filters a slice of model.Resources based on a given tag key being present on that resource.
func ResourceFilterTagKey(in []model.Resource, key string) []model.Resource {
	return FilterFunc(in, func(r model.Resource) bool {
		for _, tag := range r.Tags {
			if tag.Key == key {
				return true
			}
		}

		return false
	})
}

// ResourceFilterTagKey filters a slice of model.Resources based on a given tag key/value pair being present on that resource.
func ResourceFilterTagKeyValue(in []model.Resource, key, value string) []model.Resource {
	return FilterFunc(in, func(r model.Resource) bool {
		for _, tag := range r.Tags {
			if tag.Key == key && tag.Value == value {
				return true
			}
		}

		return false
	})
}

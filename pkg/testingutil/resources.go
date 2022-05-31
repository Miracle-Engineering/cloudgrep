package testingutil

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
)

const TestRegion = "us-east-1"
const TestTag = "test"

func AssertResourceCount(t testing.TB, resources []model.Resource, tagValue string, count int) {
	t.Helper()
	if tagValue == "" {
		resources = ResourceFilterTagKey(resources, TestTag)
	} else {
		resources = ResourceFilterTagKeyValue(resources, TestTag, tagValue)
	}

	assert.Lenf(t, resources, count, "expected %d resource(s) with tag %s=%s", count, TestTag, tagValue)
}

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

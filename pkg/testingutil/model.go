package testingutil

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
	"github.com/stretchr/testify/assert"
)

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
	return util.FilterFunc(in, func(r model.Resource) bool {
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
	return util.FilterFunc(in, func(r model.Resource) bool {
		for _, tag := range r.Tags {
			if tag.Key == key && tag.Value == value {
				return true
			}
		}

		return false
	})
}

func AssertEqualsResources(t *testing.T, a, b model.Resources) {
	assert.Equal(t, len(a), len(b))
	for _, resourceA := range a {
		resourceB := b.FindById(resourceA.Id)
		if resourceB == nil {
			t.Errorf("can't find a resource with id %v", resourceA.Id)
			return
		}
		AssertEqualsResource(t, *resourceA, *resourceB)
	}
}

func AssertEqualsResourcePter(t *testing.T, a, b *model.Resource) {
	AssertEqualsResource(t, *a, *b)
}

func AssertEqualsResource(t *testing.T, a, b model.Resource) {
	assert.Equal(t, a.Id, b.Id)
	assert.Equal(t, a.Region, b.Region)
	assert.Equal(t, a.Type, b.Type)
	jsonsEqual, err := JSONBytesEqual(a.RawData, b.RawData)
	assert.NoError(t, err)
	assert.True(t, jsonsEqual)
	assert.ElementsMatch(t, a.Tags.Clean(), b.Tags.Clean())
}

func AssertEqualsField(t *testing.T, a, b model.Field) {
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.Count, b.Count)
	assert.ElementsMatch(t, a.Values, b.Values)
}

func AssertEqualsEngineStatus(t *testing.T, expectedEngineStatus, actualEngineStatus model.EngineStatus) {
	assert.Equal(t, expectedEngineStatus.ResourceType, actualEngineStatus.ResourceType)
	assert.Equal(t, expectedEngineStatus.Status, actualEngineStatus.Status)
	assert.Equal(t, expectedEngineStatus.ErrorMessage, actualEngineStatus.ErrorMessage)
}

func AssertEqualsTag(t *testing.T, a, b *model.Tag) {
	if a == nil {
		assert.Nil(t, b)
		return
	}
	assert.Equal(t, a.Key, b.Key)
	assert.Equal(t, a.Value, b.Value)
}

func AssertEqualsTags(t *testing.T, a, b model.Tags) {
	assert.Equal(t, len(a), len(b))
	for _, tagA := range a {
		tagB := b.Find(tagA.Key)
		if tagB == nil {
			t.Errorf("can't find a tag with key %v", tagA.Key)
			return
		}
		AssertEqualsTag(t, &tagA, tagB)
	}
}

func AssertResourceFilteredCount(t testing.TB, resources []model.Resource, count int, filter ResourceFilter) []model.Resource {
	t.Helper()

	filtered := filter.Filter(resources)

	assert.Lenf(t, filtered, count, "expected %d resource(s) with filter %s", count, filter)
	return filtered
}

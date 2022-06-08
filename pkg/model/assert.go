package model

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEqualsResources(t *testing.T, a, b Resources) {
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

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func AssertEqualsResourcePter(t *testing.T, a, b *Resource) {
	AssertEqualsResource(t, *a, *b)
}

func AssertEqualsResource(t *testing.T, a, b Resource) {
	assert.Equal(t, a.Id, b.Id)
	assert.Equal(t, a.Region, b.Region)
	assert.Equal(t, a.Type, b.Type)
	jsonsEqual, err := JSONBytesEqual(a.RawData, b.RawData)
	assert.NoError(t, err)
	assert.True(t, jsonsEqual)
	assert.ElementsMatch(t, a.Tags.Clean(), b.Tags.Clean())
}

func AssertEqualsField(t *testing.T, a, b Field) {
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.Count, b.Count)
	assert.ElementsMatch(t, a.Values, b.Values)
}

func AssertEqualsResourceEvent(t *testing.T, expectedResourceEvent, actualResourceEvent ResourceEvent) {
	assert.Equal(t, expectedResourceEvent.ResourceType, actualResourceEvent.ResourceType)
	assert.Equal(t, expectedResourceEvent.FetchStatus, actualResourceEvent.FetchStatus)
	assert.Equal(t, expectedResourceEvent.ErrorMessage, actualResourceEvent.ErrorMessage)
}

func AssertEqualsEngineStatus(t *testing.T, expectedEngineStatus, actualEngineStatus EngineStatus) {
	assert.Equal(t, expectedEngineStatus.FetchStatus, actualEngineStatus.FetchStatus)
	assert.Equal(t, len(actualEngineStatus.ResourceEvents), len(expectedEngineStatus.ResourceEvents))

	for _, actualResourceEvent := range actualEngineStatus.ResourceEvents {
		var expectedResourceEvent ResourceEvent
		for _, resourceEvent := range expectedEngineStatus.ResourceEvents {
			if resourceEvent.ResourceType == actualResourceEvent.ResourceType {
				expectedResourceEvent = resourceEvent
				break
			}
		}
		assert.NotNil(t, expectedResourceEvent)
		AssertEqualsResourceEvent(t, expectedResourceEvent, actualResourceEvent)
	}
}

func AssertEqualsTag(t *testing.T, a, b *Tag) {
	if a == nil {
		assert.Nil(t, b)
		return
	}
	assert.Equal(t, a.Key, b.Key)
	assert.Equal(t, a.Value, b.Value)
}

func AssertEqualsTags(t *testing.T, a, b Tags) {
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

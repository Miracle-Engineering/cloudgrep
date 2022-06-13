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

func AssertEqualsResourceEvent(t *testing.T, re1, re2 ResourceEvent) {
	assert.Equal(t, re1.ResourceType, re2.ResourceType)
	//assert.Equal(t, re1.Provider, re2.Provider)
	assert.Equal(t, re1.FetchStatus, re2.FetchStatus)
	assert.Equal(t, re1.ErrorMessage, re2.ErrorMessage)
}

func AssertEqualsProviderStatus(t *testing.T, ps1, ps2 ProviderStatus) {
	assert.Equal(t, ps1.ProviderType, ps2.ProviderType)
	assert.Equal(t, ps1.FetchStatus, ps2.FetchStatus)
	assert.Equal(t, ps1.ErrorMessage, ps2.ErrorMessage)
	assert.Equal(t, len(ps1.ResourceEvents), len(ps2.ResourceEvents))
	for _, actualResourceEvent := range ps2.ResourceEvents {
		var expectedResourceEvent ResourceEvent
		for _, resourceEvent := range ps1.ResourceEvents {
			if resourceEvent.ResourceType == actualResourceEvent.ResourceType {
				expectedResourceEvent = resourceEvent
				break
			}
		}
		assert.NotNil(t, expectedResourceEvent)
		AssertEqualsResourceEvent(t, expectedResourceEvent, actualResourceEvent)
	}

}

func AssertEqualsEngineStatus(t *testing.T, es1, es2 EngineStatus) {
	assert.Equal(t, es1.FetchStatus, es2.FetchStatus)
	assert.Equal(t, len(es2.ProviderStatuses), len(es1.ProviderStatuses))

	for _, actualProviderStatus := range es2.ProviderStatuses {
		var expectedProviderStatus ProviderStatus
		for _, providerStatus := range es1.ProviderStatuses {
			if providerStatus.ProviderType == actualProviderStatus.ProviderType {
				expectedProviderStatus = providerStatus
				break
			}
		}
		assert.NotNil(t, expectedProviderStatus)
		AssertEqualsProviderStatus(t, expectedProviderStatus, actualProviderStatus)
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

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

func AssertEqualsEngineStatus(t *testing.T, expectedEngineStatus, actualEngineStatus EngineStatus) {
	assert.Equal(t, expectedEngineStatus.ResourceName, actualEngineStatus.ResourceName)
	assert.Equal(t, expectedEngineStatus.Status, actualEngineStatus.Status)
	assert.Equal(t, expectedEngineStatus.ErrorMessage, actualEngineStatus.ErrorMessage)
}

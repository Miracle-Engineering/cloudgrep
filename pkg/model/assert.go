package model

import (
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

func AssertEqualsResourcePter(t *testing.T, a, b *Resource) {
	AssertEqualsResource(t, *a, *b)
}

func AssertEqualsResource(t *testing.T, a, b Resource) {
	assert.Equal(t, a.Id, b.Id)
	assert.Equal(t, a.Region, b.Region)
	assert.Equal(t, a.Type, b.Type)
	assert.ElementsMatch(t, a.Properties.Clean(), b.Properties.Clean())
	assert.ElementsMatch(t, a.Tags.Clean(), b.Tags.Clean())
}

func AssertEqualsField(t *testing.T, a, b Field) {
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.Group, b.Group)
	assert.Equal(t, a.Count, b.Count)
	assert.ElementsMatch(t, a.Values, b.Values)
}

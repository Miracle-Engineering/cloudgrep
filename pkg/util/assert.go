package util

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
)

func AssertEqualsResources(t *testing.T, a, b model.Resources) {
	assert.Equal(t, len(a), len(b))
	for _, resourceA := range a {
		resourceB := b.Find(string(resourceA.Id))
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
	assert.ElementsMatch(t, a.Properties.Clean(), b.Properties.Clean())
	assert.ElementsMatch(t, a.Tags.Clean(), b.Tags.Clean())
}

func AssertEqualsTagIno(t *testing.T, a, b model.TagInfo) {
	assert.Equal(t, a.Key, b.Key)
	assert.Equal(t, a.Count, b.Count)
	assert.ElementsMatch(t, a.ResourceIds, b.ResourceIds)
	assert.ElementsMatch(t, a.Values, b.Values)
}

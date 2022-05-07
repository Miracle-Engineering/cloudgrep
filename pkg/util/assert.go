package util

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
)

func AssertEqualsResources(t *testing.T, a, b []*model.Resource) {
	assert.Equal(t, len(a), len(b))
	for i, v := range a {
		AssertEqualsResource(t, *v, *b[i])
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

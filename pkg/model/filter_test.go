package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncludeExclude(t *testing.T) {
	filter := Filter{Tags{}}
	assert.Nil(t, filter.TagsInclude())
	assert.Nil(t, filter.TagsInclude())
	filter = Filter{Tags{
		{Key: "cluster", Value: "stagging"},
		{Key: "region", Value: "us-east-1"},
	}}
	assert.Equal(t, filter.Tags, filter.TagsInclude())
	assert.Nil(t, filter.TagsExclude())
	filter = Filter{Tags{
		{Key: "cluster", Value: "stagging"},
		{Key: "region", Value: "us-east-1", Exclude: true},
	}}
	assert.Equal(t, Tags{
		{Key: "cluster", Value: "stagging"},
	}, filter.TagsInclude())
	assert.Equal(t, Tags{
		{Key: "region", Value: "us-east-1", Exclude: true},
	}, filter.TagsExclude())
}

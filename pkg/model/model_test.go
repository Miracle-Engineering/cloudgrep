package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTags(t *testing.T) {
	values := map[string]string{
		"cluster": "stagging",
		"region":  "us-east-1",
	}
	tags := newTags(values)
	assert.True(t, reflect.DeepEqual(
		Tags{
			{Key: "cluster", Value: "stagging"},
			{Key: "region", Value: "us-east-1"},
		},
		tags,
	))
	assert.Equal(t, 2, tags.DistinctKeys())

	//same tag multiple values
	values = map[string]string{
		"cluster": "stagging",
		"regions": "us-east-1,us-east2",
	}
	tags = newTags(values)
	assert.True(t, reflect.DeepEqual(
		Tags{
			{Key: "cluster", Value: "stagging"},
			{Key: "regions", Value: "us-east-1"},
			{Key: "regions", Value: "us-east2"},
		},
		tags,
	))
	assert.Equal(t, 2, tags.DistinctKeys())
}

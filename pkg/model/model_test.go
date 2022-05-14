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

func TestFindById(t *testing.T) {
	r1 := Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
	}
	r2 := Resource{
		Id: "i-124", Region: "us-east-1", Type: "test.Instance",
	}
	resources := Resources{&r1, &r2}
	assert.Equal(t, "i-123", resources.FindById("i-123").Id)
	assert.Equal(t, "i-124", resources.FindById("i-124").Id)
	assert.Nil(t, resources.FindById("i-123-not-found"))
}

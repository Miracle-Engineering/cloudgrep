package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTags(t *testing.T) {

	tags := newTags(make(map[string]string), []string{})
	assert.ElementsMatch(t, Tags{}, tags)
	assert.True(t, tags.Empty())

	//ignore empty keys
	tags = newTags(map[string]string{
		"": "stagging",
	}, []string{"", ""})
	assert.ElementsMatch(t, Tags{}, tags)
	assert.True(t, tags.Empty())

	values := map[string]string{
		"cluster": "stagging",
		"region":  "us-east-1",
	}
	tags = newTags(values, []string{})
	assert.ElementsMatch(t, Tags{
		{Key: "cluster", Value: "stagging"},
		{Key: "region", Value: "us-east-1"},
	}, tags)
	assert.Equal(t, 2, tags.DistinctKeys())

	//same tag multiple values
	values = map[string]string{
		"cluster": "stagging",
		"regions": "us-east-1,us-east2",
	}
	tags = newTags(values, []string{})
	assert.ElementsMatch(t, Tags{
		{Key: "cluster", Value: "stagging"},
		{Key: "regions", Value: "us-east-1"},
		{Key: "regions", Value: "us-east2"},
	}, tags)
	assert.Equal(t, 2, tags.DistinctKeys())

	//include exclude tags
	values = map[string]string{
		"cluster": "stagging",
	}
	excludes := []string{"team", "eks:nodegroup-name"}
	tags = newTags(values, excludes)
	assert.ElementsMatch(t, Tags{
		{Key: "cluster", Value: "stagging"},
		{Key: "team", Value: "*", Exclude: true},
		{Key: "eks:nodegroup-name", Value: "*", Exclude: true},
	}, tags)
	assert.Equal(t, 3, tags.DistinctKeys())
}

func TestFindResourceById(t *testing.T) {
	r1 := Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
	}
	r2 := Resource{
		Id: "i-124", Region: "us-east-1", Type: "test.Instance",
	}
	resources := Resources{&r1, &r2}
	assert.Equal(t, "i-123", string(resources.Find("i-123").Id))
	assert.Equal(t, "i-124", string(resources.Find("i-124").Id))
	assert.Nil(t, resources.Find("i-123-not-found"))
}

func TestFindById(t *testing.T) {
	t1 := TagInfo{
		Key: "cluster", Values: []string{"dev", "prod"},
	}
	t2 := TagInfo{
		Key: "instance_type", Values: []string{"t2.medium", "c4.xlarge"},
	}
	tagInfos := TagInfos{&t1, &t2}
	assert.Equal(t, "cluster", tagInfos.Find("cluster").Key)
	assert.Equal(t, "instance_type", tagInfos.Find("instance_type").Key)
	assert.Nil(t, tagInfos.Find("cluster-not-found"))
}

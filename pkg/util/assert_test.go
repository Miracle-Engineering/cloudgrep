package util

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
)

func TestAssertResource(t *testing.T) {
	r1 := model.Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
		Tags: []model.Tag{
			{Key: "enabled", Value: "true"},
			{Key: "eks:nodegroup", Value: "staging-default"},
		},
		Properties: []model.Property{
			{Name: "InstanceId", Value: "i-123"},
			{Name: "Architecture", Value: "x86_64"},
			{Name: "SecurityGroups[0]", Value: "sg-1"},
		},
	}
	r2 := model.Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
		Tags: []model.Tag{
			{Key: "eks:nodegroup", Value: "staging-default"},
			{Key: "enabled", Value: "true"},
		},
		Properties: []model.Property{
			{Name: "SecurityGroups[0]", Value: "sg-1"},
			{Name: "Architecture", Value: "x86_64"},
			{Name: "InstanceId", Value: "i-123"},
		},
	}
	//r1 and r2 should be equals even though the order of their tags/properties are different
	AssertEqualsResource(t, r1, r2)
}

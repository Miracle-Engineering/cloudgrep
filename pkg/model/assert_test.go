package model

import (
	"gorm.io/datatypes"
	"testing"
)

func TestAssertResource(t *testing.T) {
	r1 := Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
		Tags: []Tag{
			{Key: "enabled", Value: "true"},
			{Key: "eks:nodegroup", Value: "staging-default"},
		},
		RawData: datatypes.JSON([]byte(`{"name": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
	}
	r2 := Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
		Tags: []Tag{
			{Key: "eks:nodegroup", Value: "staging-default"},
			{Key: "enabled", Value: "true"},
		},
		RawData: datatypes.JSON([]byte(`{"name": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
	}
	//r1 and r2 should be equals even though the order of their tags/raw data are different
	AssertEqualsResource(t, r1, r2)
}

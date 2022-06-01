package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldFind(t *testing.T) {
	f1 := Field{
		Name:   "region",
		Values: nil,
	}
	f2 := Field{
		Name:   "type",
		Values: nil,
	}
	groups := FieldGroups{
		{
			Name:   "core",
			Fields: []Field{f1, f2},
		},
	}
	assert.Equal(t, "core", groups.FindGroup("core").Name)
	AssertEqualsField(t, f2, *groups.FindField("core", "type"))
}

func TestFieldsAddNullValues(t *testing.T) {
	groups := FieldGroups{
		{
			Name: "core",
			Fields: []Field{{
				Name:  "region",
				Count: 3,
				Values: []FieldValue{
					{Value: "us-east-1", Count: 2},
					{Value: "us-west-2", Count: 1},
				},
			}, {
				Name:  "type",
				Count: 3,
				Values: []FieldValue{
					{Value: "ec2.instance", Count: 3},
				},
			}, {
				Name:  "cluster",
				Count: 2,
				Values: []FieldValue{
					{Value: "dev", Count: 2},
				},
			},
			},
		},
	}
	groupsNullable := groups.AddNullValues()
	assert.Equal(t, FieldGroups{
		{
			Name: "core",
			Fields: []Field{{
				Name:  "region",
				Count: 3,
				Values: []FieldValue{
					{Value: "us-east-1", Count: 2},
					{Value: "us-west-2", Count: 1},
					//do not show (null) if all resources have this field
				},
			}, {
				Name:  "type",
				Count: 3,
				Values: []FieldValue{
					{Value: "ec2.instance", Count: 3},
				},
			}, {
				Name:  "cluster",
				Count: 2,
				Values: []FieldValue{
					{Value: "dev", Count: 2},
					//(null) count is the count of resources without this field
					{Value: "(null)", Count: 1},
				},
			},
			},
		},
	}, groupsNullable)

	// groups.
}

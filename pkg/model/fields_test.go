package model_test

import (
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/assert"
)

func TestFieldFind(t *testing.T) {
	f1 := model.Field{
		Name:   "region",
		Values: nil,
	}
	f2 := model.Field{
		Name:   "type",
		Values: nil,
	}
	groups := model.FieldGroups{
		{
			Name:   "core",
			Fields: []*model.Field{&f1, &f2},
		},
	}
	assert.Equal(t, "core", groups.FindGroup("core").Name)
	testingutil.AssertEqualsField(t, f2, *groups.FindField("core", "type"))
}

func TestFieldsAddNullValues(t *testing.T) {
	groups := model.FieldGroups{
		{
			Name: "core",
			Fields: []*model.Field{{
				Name:  "region",
				Count: 3,
				Values: []*model.FieldValue{
					{Value: "us-east-1", Count: "2"},
					{Value: "us-west-2", Count: "1"},
				},
			}, {
				Name:  "type",
				Count: 3,
				Values: []*model.FieldValue{
					{Value: "ec2.instance", Count: "3"},
				},
			}, {
				Name:  "cluster",
				Count: 2,
				Values: []*model.FieldValue{
					{Value: "dev", Count: "2"},
				},
			},
			},
		},
	}
	groupsNullable := groups.AddNullValues()
	assert.Equal(t, model.FieldGroups{
		{
			Name: "core",
			Fields: []*model.Field{{
				Name:  "region",
				Count: 3,
				Values: []*model.FieldValue{
					{Value: "us-east-1", Count: "2"},
					{Value: "us-west-2", Count: "1"},
					//do not show (missing) if all resources have this field
				},
			}, {
				Name:  "type",
				Count: 3,
				Values: []*model.FieldValue{
					{Value: "ec2.instance", Count: "3"},
				},
			}, {
				Name:  "cluster",
				Count: 2,
				Values: []*model.FieldValue{
					{Value: "dev", Count: "2"},
					//(missing) count is the count of resources without this field
					{Value: "(missing)", Count: "1"},
				},
			},
			},
		},
	}, groupsNullable)
}

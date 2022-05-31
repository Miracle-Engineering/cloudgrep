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

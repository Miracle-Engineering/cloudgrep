package model_test

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFieldFind(t *testing.T) {
	f1 := model.Field{
		Group:  "core",
		Name:   "region",
		Values: nil,
	}
	f2 := model.Field{
		Group:  "core",
		Name:   "type",
		Values: nil,
	}
	fields := model.Fields{f1, f2}
	testingutil.AssertEqualsField(t, f2, *fields.Find("core", "type"))
}

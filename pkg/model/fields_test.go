package model

import (
	"testing"
)

func TestFieldFind(t *testing.T) {
	f1 := Field{
		Group:  "core",
		Name:   "region",
		Values: nil,
	}
	f2 := Field{
		Group:  "core",
		Name:   "type",
		Values: nil,
	}
	fields := Fields{f1, f2}
	AssertEqualsField(t, f2, *fields.Find("core", "type"))
}

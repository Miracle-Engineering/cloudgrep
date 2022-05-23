package model

import (
	"testing"
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
	fields := Fields{f1, f2}
	AssertEqualsField(t, f2, *fields.Find("type"))
}

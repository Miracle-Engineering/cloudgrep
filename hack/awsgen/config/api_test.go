package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTagsAPI_Has(t *testing.T) {
	api := GetTagsAPI{}
	assert.False(t, api.Has())

	api.Call = "Foo"
	assert.True(t, api.Has())
}

package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache_Get(t *testing.T) {
	var c templateCache
	var tm *Template

	tm = c.get("foo")
	assert.Nil(t, tm)

	c.put(&Template{name: "foo"})
	tm = c.get("foo")
	assert.NotNil(t, tm)
	assert.Equal(t, "foo", tm.name)
}

func TestCache_Put_nil(t *testing.T) {
	var c templateCache
	require.Nil(t, c.templates)

	c.put(nil)
	assert.Nil(t, c.templates)

	c.put(&Template{name: "foo"})
	assert.Len(t, c.templates, 1)
}

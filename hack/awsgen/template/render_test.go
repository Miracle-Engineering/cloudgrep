package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTemplate_good(t *testing.T) {
	actual := RenderTemplate("test", 1)
	expected := "foo 1\n"

	assert.Equal(t, expected, actual)
}

func TestRenderTemplate_failedLoading(t *testing.T) {
	assert.Panics(t, func() {
		RenderTemplate("foo", nil)
	})
}

func TestRenderTemplate_failedRendering(t *testing.T) {
	assert.Panics(t, func() {
		RenderTemplate("tags.go", struct{}{})
	})
}

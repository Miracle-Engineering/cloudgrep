package template

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed templates/tags.go.tmpl
var tagsTemplate string

func TestReadTemplate_found(t *testing.T) {
	contents, err := readTemplate("tags.go")
	assert.NoError(t, err)
	assert.Equal(t, tagsTemplate, contents)
}

func TestReadTemplate_notFound(t *testing.T) {
	contents, err := readTemplate("foo.go")
	assert.ErrorContains(t, err, "cannot read template")
	assert.Equal(t, "", contents)
}

func TestParseTemplate(t *testing.T) {
	tmpl, err := parseTemplate("test", funcTestTemplate)

	assert.NoError(t, err)
	assert.NotNil(t, tmpl)
}

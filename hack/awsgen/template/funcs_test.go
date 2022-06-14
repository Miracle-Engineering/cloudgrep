package template

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/stretchr/testify/assert"
)

func TestTemplateFuncTabIndent(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		levels   int
	}{
		{
			in:       "",
			expected: "\t",
		},
		{
			in:       "foo",
			expected: "\tfoo",
		},
		{
			in:       "foo\n",
			expected: "\tfoo\n\t",
		},
		{
			in:       "foo\nbar\n",
			expected: "\t\tfoo\n\t\tbar\n\t\t",
			levels:   2,
		},
	}

	for _, test := range tests {
		if test.levels == 0 {
			test.levels = 1
		}

		name := fmt.Sprintf("l=%d;in=%s", test.levels, strconv.Quote(test.in))
		t.Run(name, func(t *testing.T) {
			actual := templateFuncTabIndent(test.levels, test.in)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestTemplateFuncQuiet(t *testing.T) {
	long := strings.Repeat("foo", 100)
	actual := templateFuncQuiet(long, long, long)
	assert.Empty(t, actual)
}

//go:embed testdata/include.tmpl
var includeTemplate string

func TestIncludeFunc(t *testing.T) {
	tmpl := includeFuncPrep(t)

	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, 2)
	assert.NoError(t, err)

	expected := strings.Repeat("foo\n", 3)
	assert.Equal(t, expected, buf.String())
}

func TestIncludeFunc_maxRecursion(t *testing.T) {
	tmpl := includeFuncPrep(t)

	err := tmpl.Execute(io.Discard, 1001)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "max recursion limit reached")
}

//go:embed testdata/funcs.tmpl
var funcTestTemplate string

func TestAddTemplateFuncs(t *testing.T) {
	tmpl := template.New("test")

	addTemplateFuncs(tmpl)

	tmpl = template.Must(tmpl.Parse(funcTestTemplate))

	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, struct{}{})
	assert.NoError(t, err)

	expected := "foo\n\n\t\tbar\n"
	assert.Equal(t, expected, buf.String())
}

func includeFuncPrep(t *testing.T) *template.Template {
	t.Helper()

	tmpl := template.New("test")

	tmpl.Funcs(sprig.TxtFuncMap())
	tmpl.Funcs(template.FuncMap{
		"include": buildIncludeFunc(tmpl),
	})

	tmpl = template.Must(tmpl.Parse(includeTemplate))

	return tmpl
}

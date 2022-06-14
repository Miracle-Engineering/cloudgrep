package template

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func addTemplateFuncs(t *template.Template) {
	t.Funcs(sprig.TxtFuncMap())
	t.Funcs(template.FuncMap{
		"include":   buildIncludeFunc(t),
		"quiet":     templateFuncQuiet,
		"tabindent": templateFuncTabIndent,
	})
}

// recursionMaxNums is the maximum allowed recursion level for the "include" template function.
const recursionMaxNums = 1000

// buildIncludeFunc generates the "include" template function implementation.
// It returns the function because the function depends on the template.Template value to be created.
func buildIncludeFunc(t *template.Template) any {
	includedNames := make(map[string]int)

	return func(name string, data any) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", fmt.Errorf("max recursion limit reached for %s", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := t.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}
}

// templateFuncTabIndent (or "tabindent" within templates) indents each line in the input
// (including the first and last, even if empty) a certain number of levels using tab characters.
func templateFuncTabIndent(levels int, v string) string {
	pad := strings.Repeat("\t", levels)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

// templateFuncQuiet (or "quiet" within templates) acts as a /dev/null for the input.
// Useful for silencing the output from a variable assignment.
func templateFuncQuiet(_ ...any) string {
	return ""
}

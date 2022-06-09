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

const recursionMaxNums = 1000

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

func templateFuncTabIndent(levels int, v string) string {
	pad := strings.Repeat("\t", levels)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func templateFuncQuiet(_ ...any) string {
	return ""
}

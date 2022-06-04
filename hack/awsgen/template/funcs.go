package template

import (
	"fmt"
	"strings"
	"text/template"
)

const recursionMaxNums = 1000

func buildIncludeFunc(t *template.Template) any {
	includedNames := make(map[string]int)

	return func(name string, data any) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", fmt.Errorf("rendering template has a nested reference name: %s", name)
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

package template

import (
	"fmt"
	"path"
	"text/template"

	"github.com/Masterminds/sprig"
)

func fetchTemplate(name string) *template.Template {
	if tmpl, ok := templates[name]; ok {
		return tmpl
	}

	t := loadTemplate(name)

	templates[name] = t

	return t
}

func loadTemplate(name string) *template.Template {
	templatePath := path.Join("templates", name+".tmpl")
	templateContents, err := content.ReadFile(templatePath)
	if err != nil {
		panic(fmt.Errorf("cannot read template %s: %w", name, err))
	}

	t := template.New(name)

	addFuncs(t)

	_, err = t.Parse(string(templateContents))
	if err != nil {
		panic(fmt.Errorf("error parsing template %s: %w", name, err))
	}

	t.Option("missingkey=error")

	return t
}

func addFuncs(t *template.Template) {
	t.Funcs(sprig.TxtFuncMap())
	t.Funcs(template.FuncMap{
		"include":   buildIncludeFunc(t),
		"tabindent": templateFuncTabIndent,
	})
}

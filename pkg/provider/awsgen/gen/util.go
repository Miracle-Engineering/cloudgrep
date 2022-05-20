package main

import (
	"bytes"
	"embed"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed *
var content embed.FS

func renderTemplate(name string, config interface{}, funcs ...template.FuncMap) string {
	tmpl := fetchTemplate(name, funcs...)
	buf := &bytes.Buffer{}

	err := tmpl.Execute(buf, config)
	if err != nil {
		panic(fmt.Errorf("cannot render template %s: %w", config, err))
	}

	return buf.String()
}

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
}

func fetchTemplate(name string, funcs ...template.FuncMap) *template.Template {
	if tmpl, ok := templates[name]; ok {
		return tmpl
	}

	templatePath := path.Join("templates", name+".tmpl")
	templateContents, err := content.ReadFile(templatePath)
	if err != nil {
		panic(fmt.Errorf("cannot read template %s: %w", name, err))
	}

	t := template.New(name)
	for _, f := range funcs {
		t.Funcs(f)
	}

	t.Funcs(sprig.TxtFuncMap())
	t.Funcs(template.FuncMap{
		"include": buildIncludeFunc(t),
	})

	_, err = t.Parse(string(templateContents))
	if err != nil {
		panic(fmt.Errorf("error parsing template %s: %w", name, err))
	}

	t.Option("missingkey=error")

	templates[name] = t

	return t
}

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

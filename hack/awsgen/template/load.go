package template

import (
	"embed"
	"fmt"
	"path"
	"text/template"
)

//go:embed templates
var content embed.FS

func readTemplate(name string) (string, error) {
	templatePath := path.Join("templates", name+".tmpl")
	contents, err := content.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("cannot read template %s: %w", name, err)
	}

	return string(contents), nil
}

func parseTemplate(name string, contents string) (*template.Template, error) {
	t := template.New(name)
	addTemplateFuncs(t)

	_, err := t.Parse(contents)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %s: %w", name, err)
	}

	t.Option("missingkey=error")

	return t, nil
}

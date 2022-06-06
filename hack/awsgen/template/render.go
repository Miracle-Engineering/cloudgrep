package template

import (
	"bytes"
	"embed"
	"fmt"
)

//go:embed templates
var content embed.FS

func RenderTemplate(name string, config interface{}) string {
	tmpl := getTemplate(name)
	buf := &bytes.Buffer{}

	err := tmpl.Execute(buf, config)
	if err != nil {
		panic(fmt.Errorf("cannot render template %s: %w", config, err))
	}

	return buf.String()
}

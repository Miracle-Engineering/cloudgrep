package template

import (
	"bytes"
	"fmt"
)

func RenderTemplate(name string, config interface{}) string {
	tmpl := fetchTemplate(name)
	buf := &bytes.Buffer{}

	err := tmpl.Execute(buf, config)
	if err != nil {
		panic(fmt.Errorf("cannot render template %s: %w", config, err))
	}

	return buf.String()
}

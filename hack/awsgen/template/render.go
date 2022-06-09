package template

import (
	"fmt"
)

func RenderTemplate(name string, data any) string {
	tmpl, err := NewTemplate(name)
	if err != nil {
		panic(fmt.Errorf("error loading template %s: %w", name, err))
	}

	result, err := tmpl.Render(data)
	if err != nil {
		panic(fmt.Errorf("error rendering template %s: %w", name, err))
	}

	return result
}

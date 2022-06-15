package template

import (
	"bytes"
	"fmt"
	"text/template"
)

type Template struct {
	name string
	t    *template.Template
}

func NewTemplate(name string) (*Template, error) {
	if cached := cache.get(name); cached != nil {
		return cached, nil
	}

	t, err := newTemplateNoCache(name)
	if err == nil && t != nil {
		cache.put(t)
	}

	return t, err
}

func newTemplateNoCache(name string) (*Template, error) {
	contents, err := readTemplate(name)
	if err != nil {
		return nil, fmt.Errorf("error creating template: %w", err)
	}

	parsed, err := parseTemplate(name, contents)
	if err != nil {
		return nil, fmt.Errorf("error creating template: %w", err)
	}

	if parsed == nil {
		panic("unexpected nil parsed")
	}

	return &Template{
		name: name,
		t:    parsed,
	}, nil
}

func (t *Template) Render(data any) (string, error) {
	buf := &bytes.Buffer{}

	err := t.t.Execute(buf, data)
	if err != nil {
		return "", fmt.Errorf("cannot render template %s: %w", t.name, err)
	}

	return buf.String(), nil
}

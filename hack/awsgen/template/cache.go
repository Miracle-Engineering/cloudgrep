package template

import "text/template"

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
}

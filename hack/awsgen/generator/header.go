package generator

import "github.com/run-x/cloudgrep/hack/awsgen/template"

func (g *Generator) generateFileHeader(pkg string, imports []Import) string {
	data := struct {
		Package string
		Imports []Import
	}{
		Package: pkg,
		Imports: imports,
	}

	return template.RenderTemplate("header.go", data)
}

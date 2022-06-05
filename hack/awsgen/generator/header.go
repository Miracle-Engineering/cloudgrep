package generator

import "github.com/run-x/cloudgrep/hack/awsgen/template"

type fileHeader struct {
	Package string
	Imports []Import
}

func (g *Generator) generateFileHeader(header fileHeader) string {
	return template.RenderTemplate("header.go", header)
}

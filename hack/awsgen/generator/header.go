package generator

import (
	"github.com/run-x/cloudgrep/hack/awsgen/template"
	"github.com/run-x/cloudgrep/hack/awsgen/util"
)

func (g Generator) generateFileHeader(pkg string, imports []util.Import) string {
	data := struct {
		Package string
		Imports [][]util.Import
	}{
		Package: pkg,
	}

	util.SortImports(imports)
	data.Imports = util.GroupImports(imports).Groups()

	return template.RenderTemplate("header.go", data)
}

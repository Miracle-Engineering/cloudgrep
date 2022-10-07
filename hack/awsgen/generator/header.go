package generator

import (
	"github.com/juandiegopalomino/cloudgrep/hack/awsgen/template"
	"github.com/juandiegopalomino/cloudgrep/hack/awsgen/util"
)

// generateFileHeader returns the package statement and import statement for a single file
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

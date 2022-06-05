package generator

import (
	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
)

func (g *Generator) generateMainFile(services []config.ServiceConfig) string {
	c := struct {
		Services []config.ServiceConfig
	}{
		Services: services,
	}

	header := g.generateFileHeader(fileHeader{
		Package: "aws",
		Imports: []Import{
			// {
			// 	Path: "reflect",
			// },
			// {
			// 	Path: "github.com/run-x/cloudgrep/pkg/provider/awsgen/registry",
			// },
		},
	})

	body := template.RenderTemplate("all.go", c)

	return header + "\n" + body
}

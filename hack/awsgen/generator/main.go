package generator

import (
	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
)

func (g *Generator) generateMainFile(services []config.ServiceConfig) string {
	c := struct {
		Package  string
		Services []config.ServiceConfig
	}{
		Package:  "generated",
		Services: services,
	}

	return template.RenderTemplate("all.go", c)
}

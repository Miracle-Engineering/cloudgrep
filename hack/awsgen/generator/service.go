package generator

import (
	"fmt"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
)

func (g Generator) generateServiceHeader(svc config.ServiceConfig) string {
	importPath := "github.com/aws/aws-sdk-go-v2/service/%s"

	header := g.generateFileHeader(fileHeader{
		Package: "aws",
		Imports: simpleImports([]string{
			fmt.Sprintf(importPath, svc.Name),
			"context",
			"fmt",

			"github.com/run-x/cloudgrep/pkg/resourceconverter",
			"github.com/run-x/cloudgrep/pkg/model",
		}),
	})

	c := serviceHeaderConfig{
		Service: svc.Name,
	}

	for _, t := range svc.Types {
		c.Types = append(c.Types, serviceHeaderTypeConfig{
			Name:     svc.Name + "." + t.Name,
			FuncName: fetchFuncName(svc, t),
			Type:     t.Name,
			ID:       t.ListAPI.IDField,
			Tags:     t.Tags,
			Global:   t.Global,
		})
	}

	body := template.RenderTemplate("service.go", c)

	return header + "\n" + body
}

type serviceHeaderConfig struct {
	Service string
	Types   []serviceHeaderTypeConfig
}

type serviceHeaderTypeConfig struct {
	Name     string
	FuncName string
	Type     string
	ID       string
	Tags     config.TagField
	Global   bool
}

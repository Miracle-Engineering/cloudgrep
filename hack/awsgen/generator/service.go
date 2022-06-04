package generator

import (
	"fmt"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
)

func (g Generator) generateServiceHeader(svc config.ServiceConfig) string {
	importPath := "github.com/aws/aws-sdk-go-v2/service/%s/types"
	c := serviceHeaderConfig{
		Package: "generated",
		Service: svc.Name,
	}
	c.Imports = []Import{
		{
			Path: fmt.Sprintf(importPath, svc.Name),
		},
	}

	for _, t := range svc.Types {
		c.Types = append(c.Types, serviceHeaderTypeConfig{
			Name: svc.Name + "." + t.Name,
			Type: t.Name,
			ID:   t.ListAPI.IDField,
			Tag:  t.ListAPI.TagField,
		})
	}

	return template.RenderTemplate("service.go", c)
}

type serviceHeaderConfig struct {
	Package string
	Imports []Import
	Service string
	Types   []serviceHeaderTypeConfig
}

type serviceHeaderTypeConfig struct {
	Name string
	Type string
	ID   string
	Tag  string
}

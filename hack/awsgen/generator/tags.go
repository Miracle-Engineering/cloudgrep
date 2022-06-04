package generator

import (
	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
)

type typeTagFunctionConfig struct {
	Name         string
	Action       string
	Description  string
	Type         string
	Client       string
	Service      string
	OutputKey    string
	InputIDField string
	IDField      string
}

func (g Generator) generateTypeTagFunction(svc config.ServiceConfig, t config.TypeConfig) string {
	c := typeTagFunctionConfig{
		Name:         "Fetch" + t.Name + "Tags",
		Action:       t.GetTagsAPI.Call,
		Description:  t.Description,
		Type:         t.Name,
		Client:       svc.Name + "Client",
		Service:      svc.Name,
		OutputKey:    t.GetTagsAPI.OutputKey,
		InputIDField: t.GetTagsAPI.InputIDField,
		IDField:      t.ListAPI.IDField,
	}

	return template.RenderTemplate("tags.go", c)
}

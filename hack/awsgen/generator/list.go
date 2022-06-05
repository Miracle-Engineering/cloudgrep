package generator

import (
	"github.com/run-x/cloudgrep/hack/awsgen/template"
	"github.com/run-x/cloudgrep/hack/awsgen/util"
)

type listFuncConfig struct {
	ResourceName string
	FuncName     string
	Action       string
	Paginated    bool
	Description  string
	Type         string
	Client       string
	ServicePkg   string
	ProviderName string
	OutputKey    util.RecursiveAppend
}

func (g Generator) generateListFunction(config listFuncConfig) string {
	name := "list.go"
	if config.Paginated {
		name = "list-paginated.go"
	}

	return template.RenderTemplate(name, config)
}

package generator

import (
	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
)

func (g Generator) generateRegisterFile(services []config.ServiceConfig) string {
	data := struct {
		ProviderName      string
		RegisterFuncNames []string
	}{
		ProviderName: "Provider",
	}

	for _, service := range services {
		data.RegisterFuncNames = append(data.RegisterFuncNames, registerFuncName(service))
	}

	return template.RenderTemplate("register.go", data)
}

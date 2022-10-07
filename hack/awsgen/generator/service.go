package generator

import (
	"strings"

	"github.com/juandiegopalomino/cloudgrep/hack/awsgen/config"
	"github.com/juandiegopalomino/cloudgrep/hack/awsgen/template"
	"github.com/juandiegopalomino/cloudgrep/hack/awsgen/util"
)

// generateService generates the file for a specific resource
func (g Generator) generateService(service config.Service) string {
	var imports util.ImportSet

	buf := &strings.Builder{}

	reg, regImports := g.generateServiceRegister(service)
	buf.WriteString(reg)
	imports.Merge(regImports)

	for _, t := range service.Types {
		f, typeImports := g.generateType(service, t)
		buf.WriteString(f)
		imports.Merge(typeImports)
	}

	header := g.generateFileHeader(PackageName, imports.Get())

	return header + "\n" + buf.String()
}

// generateServiceRegister generates the services "registration" function that returns data on each type defined in a specific service.
func (g Generator) generateServiceRegister(service config.Service) (string, util.ImportSet) {
	data := struct {
		ProviderName string
		ServiceName  string
		EndpointID   string
		FuncName     string

		Types []typeRegisterInfo
	}{
		ProviderName: ProviderStructName,
		EndpointID:   service.EndpointID,
		ServiceName:  service.Name,
		FuncName:     registerFuncName(service),
	}

	var imports util.ImportSet

	for _, typ := range service.Types {
		global := service.Global
		if typ.Global != nil {
			global = *typ.Global
		}

		data.Types = append(data.Types, typeRegisterInfo{
			ResourceName:   resourceName(service, typ),
			FetchFuncName:  fetchFuncName(service, typ),
			IDField:        typ.ListAPI.IDField,
			DisplayIDField: typ.ListAPI.DisplayIDField,
			Global:         global,
			Tags:           typ.ListAPI.Tags,
		})

		if !typ.GetTagsAPI.Tags.Zero() {
			imports.AddPath("github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter")
		}
	}

	return template.RenderTemplate("service-register.go", data), imports
}

type typeRegisterInfo struct {
	ResourceName  string
	FetchFuncName string

	IDField        config.Field
	DisplayIDField config.Field
	Global         bool
	Tags           *config.TagField
}

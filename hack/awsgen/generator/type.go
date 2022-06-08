package generator

import (
	"strings"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
	"github.com/run-x/cloudgrep/hack/awsgen/util"
)

func (g Generator) generateType(service config.Service, typ config.Type) (string, util.ImportSet) {
	var imports util.ImportSet

	buf := &strings.Builder{}

	listFunc, listImports := g.generateTypeListFunction(service, typ)
	buf.WriteString(listFunc)
	imports.Merge(listImports)

	tagFunc, tagImports := g.generateTypeTagFunction(service, typ)
	buf.WriteString(tagFunc)
	imports.Merge(tagImports)

	return buf.String(), imports
}

func (g Generator) generateTypeListFunction(service config.Service, typ config.Type) (string, util.ImportSet) {
	data := struct {
		ResourceName string
		Description  string

		FuncName     string
		ProviderName string

		ServicePkg string
		APIAction  string
		Paginated  bool

		OutputKey   *util.RecursiveAppend
		TagFuncName string
	}{
		ResourceName: resourceName(service, typ),
		Description:  typ.Description,

		FuncName:     fetchFuncName(service, typ),
		ProviderName: ProviderStructName,

		ServicePkg: service.ServicePackage,
		APIAction:  typ.ListAPI.Call,
		Paginated:  typ.ListAPI.Pagination,
		OutputKey: &util.RecursiveAppend{
			Keys: typ.ListAPI.OutputKey,
		},
	}

	var imports util.ImportSet
	imports.AddPath("context")
	imports.AddPath("fmt")
	imports.AddPath(awsServicePackage(service.ServicePackage))
	imports.AddPath("github.com/run-x/cloudgrep/pkg/resourceconverter")
	imports.AddPath("github.com/run-x/cloudgrep/pkg/model")

	if typ.GetTagsAPI.Has() {
		data.TagFuncName = tagFuncName(service, typ)
	}

	return template.RenderTemplate("list.go", data), imports
}

func (g Generator) generateTypeTagFunction(service config.Service, typ config.Type) (string, util.ImportSet) {
	if !typ.GetTagsAPI.Has() {
		return "", nil
	}

	if typ.GetTagsAPI.TagField == nil {
		panic("unexpected nil getTagsApi.tags")
	}

	data := struct {
		ResourceName string
		Description  string

		FuncName     string
		ProviderName string

		ServicePkg           string
		APIAction            string
		SDKType              string
		AllowedAPIErrorCodes []string

		InputIDField    config.Field
		ResourceIDField config.Field
		TagField        config.TagField
	}{
		ResourceName: resourceName(service, typ),
		Description:  typ.Description,

		FuncName:     tagFuncName(service, typ),
		ProviderName: ProviderStructName,

		ServicePkg:           service.ServicePackage,
		APIAction:            typ.GetTagsAPI.Call,
		SDKType:              typ.GetTagsAPI.ResourceType,
		AllowedAPIErrorCodes: typ.GetTagsAPI.AllowedAPIErrorCodes,

		InputIDField:    typ.GetTagsAPI.InputIDField,
		ResourceIDField: typ.ListAPI.IDField,
		TagField:        *typ.GetTagsAPI.TagField,
	}

	var imports util.ImportSet
	imports.AddPath("context")
	imports.AddPath("fmt")
	imports.AddPath(awsServicePackage(service.ServicePackage))
	imports.AddPath(awsServicePackage(service.ServicePackage, "types"))
	imports.AddPath("github.com/run-x/cloudgrep/pkg/model")

	if len(typ.GetTagsAPI.AllowedAPIErrorCodes) > 0 {
		imports.AddPath("github.com/aws/smithy-go")
		imports.AddPath("errors")
	}

	return template.RenderTemplate("tags.go", data), imports
}
package generator

import (
	"strings"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
	"github.com/run-x/cloudgrep/hack/awsgen/util"
)

// generateType generates the functions for a specific type
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

// generateTypeListFunction generates the code for listing a specific type
func (g Generator) generateTypeListFunction(service config.Service, typ config.Type) (string, util.ImportSet) {
	data := struct {
		SDKType      string
		ResourceName string

		FuncName     string
		ProviderName string

		ServicePkg     string
		APIAction      string
		Paginated      bool
		InputOverrides config.InputOverrides

		OutputKey   *util.RecursiveAppend[config.Field]
		TagFuncName string
	}{
		SDKType:      sdkType(typ),
		ResourceName: resourceName(service, typ),

		FuncName:     fetchFuncName(service, typ),
		ProviderName: ProviderStructName,

		ServicePkg:     service.ServicePackage,
		APIAction:      typ.ListAPI.Call,
		Paginated:      typ.ListAPI.Pagination,
		InputOverrides: typ.ListAPI.InputOverrides,

		OutputKey: &util.RecursiveAppend[config.Field]{
			Keys: typ.ListAPI.OutputKey,
		},
	}

	var imports util.ImportSet
	imports.AddPath("context")
	imports.AddPath("fmt")
	imports.AddPath(awsServicePackage(service.ServicePackage))
	imports.AddPath(awsServicePackage(service.ServicePackage, "types"))
	imports.AddPath("github.com/run-x/cloudgrep/pkg/resourceconverter")
	imports.AddPath("github.com/run-x/cloudgrep/pkg/model")

	if typ.GetTagsAPI.Has() {
		data.TagFuncName = tagFuncName(service, typ)
	}

	return template.RenderTemplate("list.go", data), imports
}

// generateTypeTagFunction generates the code for fetching tags for a specific type
func (g Generator) generateTypeTagFunction(service config.Service, typ config.Type) (string, util.ImportSet) {
	if !typ.GetTagsAPI.Has() {
		return "", nil
	}

	if typ.GetTagsAPI.Tags == nil {
		panic("unexpected nil getTagsApi.tags")
	}

	data := struct {
		ResourceName string

		FuncName     string
		ProviderName string

		ServicePkg           string
		APIAction            string
		SDKType              string
		AllowedAPIErrorCodes []string
		InputOverrides       config.InputOverrides

		InputIDField    config.Field
		ResourceIDField config.Field
		Tags            config.TagField
	}{
		ResourceName: resourceName(service, typ),

		FuncName:     tagFuncName(service, typ),
		ProviderName: ProviderStructName,

		ServicePkg:           service.ServicePackage,
		APIAction:            typ.GetTagsAPI.Call,
		SDKType:              sdkType(typ),
		AllowedAPIErrorCodes: typ.GetTagsAPI.AllowedAPIErrorCodes,
		InputOverrides:       typ.GetTagsAPI.InputOverrides,

		InputIDField:    typ.GetTagsAPI.InputIDField,
		ResourceIDField: typ.ListAPI.IDField,
		Tags:            *typ.GetTagsAPI.Tags,
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

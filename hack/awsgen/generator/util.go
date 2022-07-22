package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
)

const PackageName = "aws"
const ProviderStructName = "Provider"

// linenumbers adds line number comments to the start of each line in `in`.
// Uses go's multiline comment format.
func linenumbers(in string) string {
	b := strings.Builder{}
	lines := strings.Split(in, "\n")
	chars := len(strconv.Itoa(len(lines)))

	for idx, line := range lines {
		lineNum := idx + 1
		lineText := strconv.Itoa(lineNum)
		paddingNeeded := chars - len(lineText)
		padding := strings.Repeat(" ", paddingNeeded)
		b.WriteString("/* ")
		b.WriteString(padding)
		b.WriteString(lineText)
		b.WriteString(" */ ")
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

// fetchFuncName returns the name of the fetch (aka list) functions for a specific type
func fetchFuncName(svc config.Service, typ config.Type) string {
	return lowerCamelCaseJoin(
		"fetch",
		svc.Name,
		typ.Name,
	)
}

// tagFuncName returns the name of the tags function for a specific type
func tagFuncName(svc config.Service, typ config.Type) string {
	return lowerCamelCaseJoin(
		"getTags",
		svc.Name,
		typ.Name,
	)
}

// awsServicePackage returns the package import string for a specific service's package in the AWS SDK,
// with support for optional subpackages (such as "types")
func awsServicePackage(service string, subPackages ...string) string {
	pkg := "github.com/aws/aws-sdk-go-v2/service/" + service

	for _, subPackage := range subPackages {
		pkg += "/" + subPackage
	}

	return pkg
}

// resourceName returns Cloudgrep's identifier for a resource type, which includes the service
func resourceName(service config.Service, typ config.Type) string {
	return fmt.Sprintf("%s.%s", service.Name, typ.Name)
}

// registerFuncName returns the name of the Go func that returns type mapping data for a specific service.
func registerFuncName(svc config.Service) string {
	return lowerCamelCaseJoin(
		"register",
		svc.Name,
	)
}

// only ASCII supported
func lowerCamelCaseJoin(parts ...string) string {
	var out []string
	for idx, part := range parts {
		if len(part) == 0 {
			continue
		}

		if idx > 0 {
			firstChar := part[0:1]
			firstChar = strings.ToUpper(firstChar)
			part = firstChar + part[1:]
		}

		out = append(out, part)
	}

	return strings.Join(out, "")
}

func sdkType(typ config.Type) string {
	if typ.ListAPI.SDKType != "" {
		return typ.ListAPI.SDKType
	}

	return typ.Name
}

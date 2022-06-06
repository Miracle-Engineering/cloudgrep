package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
)

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

func fetchFuncName(svc config.ServiceConfig, typ config.TypeConfig) string {
	return fmt.Sprintf(
		"fetch_%s_%s",
		svc.Name,
		typ.Name,
	)
}

func tagFuncName(svc config.ServiceConfig, typ config.TypeConfig) string {
	return fmt.Sprintf(
		"getTags_%s_%s",
		svc.Name,
		typ.Name,
	)
}

func awsServicePackage(service string, subPackages ...string) string {
	pkg := "github.com/aws/aws-sdk-go-v2/service/" + service

	for _, subPackage := range subPackages {
		pkg += "/" + subPackage
	}

	return pkg
}

func resourceName(service config.ServiceConfig, typ config.TypeConfig) string {
	return fmt.Sprintf("%s.%s", service.Name, typ.Name)
}

func registerFuncName(svc config.ServiceConfig) string {
	return "register_" + svc.Name
}

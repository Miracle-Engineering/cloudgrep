package main

import (
	"fmt"
	"go/format"
	"strconv"
	"strings"
)

func main() {
	services := loadServices()

	for _, service := range services {
		out := generateService(service)

		outRaw, err := format.Source([]byte(out))
		if err != nil {
			fmt.Println(linenumbers(out))
			panic(fmt.Errorf("error formatting service %s: %w", service.Name, err))
		}
		out = string(outRaw)

		writeFile(service.Name, out)
	}

	writeFile("services", generateMainFile(services))

}

func writeFile(name, contents string) {
	name = "zz_" + name + ".go"
	fmt.Printf("// %s\n", name)
	fmt.Println(contents)
	fmt.Println()
}

func linenumbers(in string) string {
	b := strings.Builder{}
	lines := strings.Split(in, "\n")
	chars := len(strconv.Itoa(len(lines)))

	for idx, line := range lines {
		lineNum := idx + 1
		lineText := strconv.Itoa(lineNum)
		paddingNeeded := chars - len(lineText)
		padding := strings.Repeat(" ", paddingNeeded)
		b.WriteString(padding)
		b.WriteString(lineText)
		b.WriteString(" ")
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

func generateService(service ServiceConfig) string {
	buf := &strings.Builder{}

	buf.WriteString(generateServiceHeader(service))

	for _, t := range service.Types {
		c := listFuncConfig{
			Name:        "Fetch" + t.Name,
			Action:      t.ListAPI.Call,
			Paginated:   t.ListAPI.Pagination,
			Description: t.Description,
			Type:        t.Name,
			Client:      service.Name + "Client",
			Service:     service.Name,
			OutputKey: RecursiveAppend{
				Keys: t.ListAPI.OutputKey,
				Root: "output",
			},
		}

		f := generateListFunction(c)
		buf.WriteString(f)

		if t.GetTagsAPI.Call != "" {
			f = generateTypeTagFunction(service, t)
			buf.WriteString(f)
		}
	}

	return buf.String()
}

type listFuncConfig struct {
	Name        string
	Action      string
	Paginated   bool
	Description string
	Type        string
	Client      string
	Service     string
	OutputKey   RecursiveAppend
}

type RecursiveAppend struct {
	Idx  int
	Keys []string
	Root string
}

func (r RecursiveAppend) IsLast() bool {
	return r.Idx >= len(r.Keys)-1
}

func (r RecursiveAppend) IterVar() string {
	if r.Idx > 0 {
		return r.varFor(r.Idx - 1)
	}

	return r.Root
}

func (r RecursiveAppend) Current() string {
	return r.Keys[r.Idx]
}

func (r RecursiveAppend) NextIterVar() string {
	return r.varFor(r.Idx)
}

func (r RecursiveAppend) varFor(i int) string {
	return "item_" + strconv.Itoa(i)
}

func (r RecursiveAppend) Next() (RecursiveAppend, error) {
	if r.IsLast() {
		return r, fmt.Errorf("end of recursive append keys")
	}

	r.Idx++
	return r, nil
}

func (r RecursiveAppend) WithRoot(root string) RecursiveAppend {
	r.Root = root
	return r
}

func generateListFunction(config listFuncConfig) string {
	name := "list.go"
	if config.Paginated {
		name = "list-paginated.go"
	}

	return renderTemplate(name, config, map[string]any{
		"tabindent": templateFuncTabIndent,
	})
}

func templateFuncTabIndent(levels int, v string) string {
	pad := strings.Repeat("\t", levels)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

type serviceHeaderConfig struct {
	Package string
	Imports []Import
	Service string
	Types   []ServiceHeaderTypeConfig
}

type Import struct {
	Path string
	As   string
}

type ServiceHeaderTypeConfig struct {
	Name string
	Type string
	ID   string
	Tag  string
}

func generateServiceHeader(svc ServiceConfig) string {
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
		c.Types = append(c.Types, ServiceHeaderTypeConfig{
			Name: svc.Name + "." + t.Name,
			Type: t.Name,
			ID:   t.ListAPI.IDField,
			Tag:  t.ListAPI.TagField,
		})
	}

	return renderTemplate("service.go", c)
}

func generateMainFile(services map[string]ServiceConfig) string {
	c := struct {
		Package  string
		Services []string
	}{
		Package:  "generated",
		Services: MapKeys(services),
	}

	return renderTemplate("all.go", c)
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

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

func generateTypeTagFunction(svc ServiceConfig, t TypeConfig) string {
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

	return renderTemplate("tags.go", c)
}

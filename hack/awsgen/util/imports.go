package util

import (
	"sort"
	"strings"
)

type Import struct {
	Path string
	As   string
}

func SimpleImports(paths []string) []Import {
	imports := make([]Import, 0, len(paths))
	for _, path := range paths {
		imports = append(imports, Import{
			Path: path,
		})
	}

	return imports
}

type ImportSet map[Import]struct{}

func (i *ImportSet) Add(newImport Import) {
	if *i == nil {
		*i = make(map[Import]struct{})
	}

	(*i)[newImport] = struct{}{}
}

func (i *ImportSet) AddPath(path string) {
	i.Add(Import{Path: path})
}

func (i ImportSet) Get() []Import {
	out := make([]Import, 0, len(i))
	for imp := range i {
		out = append(out, imp)
	}

	return out
}

func (i *ImportSet) Merge(other ImportSet) {
	for imp := range other {
		i.Add(imp)
	}
}

type GroupedImports struct {
	StandardLib []Import
	ThirdParty  []Import
	Module      []Import
}

const modulePrefix = "github.com/run-x/cloudgrep"

func GroupImports(imports []Import) GroupedImports {
	var grouped GroupedImports

	for _, imp := range imports {
		if strings.HasPrefix(imp.Path, modulePrefix) {
			grouped.Module = append(grouped.Module, imp)
		} else if strings.Contains(imp.Path, ".") {
			grouped.ThirdParty = append(grouped.ThirdParty, imp)
		} else {
			grouped.StandardLib = append(grouped.StandardLib, imp)
		}
	}

	return grouped
}

func (i GroupedImports) Len() int {
	var count int
	for _, group := range i.Groups() {
		count += len(group)
	}

	return count
}

func (i GroupedImports) Groups() [][]Import {
	return [][]Import{
		i.StandardLib,
		i.ThirdParty,
		i.Module,
	}
}

func SortImports(imports []Import) {
	less := func(i, j int) bool {
		return strings.Compare(imports[i].Path, imports[j].Path) < 0
	}

	sort.SliceStable(imports, less)
}

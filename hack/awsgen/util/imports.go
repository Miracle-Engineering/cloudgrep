package util

import (
	"sort"
	"strings"
)

// Import represents an import of a package in a Go file
type Import struct {
	Path string
	As   string
}

// ImportSet is a set of Imports (for removing duplicates).
// Its zero value can be used without initialization.
type ImportSet map[Import]struct{}

// Add adds the given import to the set, ignoring duplicates.
func (i *ImportSet) Add(newImport Import) {
	if *i == nil {
		*i = make(map[Import]struct{})
	}

	(*i)[newImport] = struct{}{}
}

// AddPath adds the given import path to the set.
func (i *ImportSet) AddPath(path string) {
	i.Add(Import{Path: path})
}

// Get returns a new slice of all the imports stored in the set.
// The order is not specified; it is recommended to call SortImports on the returned value.
func (i ImportSet) Get() []Import {
	out := make([]Import, 0, len(i))
	for imp := range i {
		out = append(out, imp)
	}

	return out
}

// Merge adds all Import values from `other` into this set.
func (i *ImportSet) Merge(other ImportSet) {
	for imp := range other {
		i.Add(imp)
	}
}

// GroupedImports groups imports into different sections, for cleaner import blocks.
type GroupedImports struct {
	StandardLib []Import
	ThirdParty  []Import
	Module      []Import
}

// modulePrefix holds the prefix for packages that are considered "in this module".
// If this module is renamed, this value must be updated.
const modulePrefix = "github.com/juandiegopalomino/cloudgrep/"

// GroupImports groups the passed imports into a GroupedImports.
// Each group maintains the relative order of the imports.
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

// Len returns the total number of imports
func (i GroupedImports) Len() int {
	var count int
	for _, group := range i.Groups() {
		count += len(group)
	}

	return count
}

// Groups returns a slice of slice of Import, with each top-level slice being a different group.
// The groups are ordered as standard library, third party, and then module imports.
func (i GroupedImports) Groups() [][]Import {
	return [][]Import{
		i.StandardLib,
		i.ThirdParty,
		i.Module,
	}
}

// SortImports performs an in-place sort on the passed imports, sorting each by the Path.
func SortImports(imports []Import) {
	less := func(i, j int) bool {
		return strings.Compare(imports[i].Path, imports[j].Path) < 0
	}

	sort.SliceStable(imports, less)
}

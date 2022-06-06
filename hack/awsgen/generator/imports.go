package generator

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

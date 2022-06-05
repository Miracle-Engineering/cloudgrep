package generator

import (
	"strconv"
	"strings"
)

type Import struct {
	Path string
	As   string
}

func simpleImports(raw []string) []Import {
	imports := make([]Import, 0, len(raw))
	for _, path := range raw {
		imports = append(imports, Import{
			Path: path,
		})
	}

	return imports
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
		b.WriteString("/* ")
		b.WriteString(padding)
		b.WriteString(lineText)
		b.WriteString(" */ ")
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

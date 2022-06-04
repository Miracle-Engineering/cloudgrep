package writer

import (
	"fmt"
	"io"
)

type streamWriter struct {
	dest io.Writer
}

func NewStreamWriter(dest io.Writer) Writer {
	if dest == nil {
		dest = io.Discard
	}

	return &streamWriter{dest: dest}
}

func (w *streamWriter) WriteFile(name string, contents []byte) error {
	header := fmt.Sprintf("// %s\n", fileName(name))
	footer := "\n\n"

	return multiWrite(w.dest, []byte(header), contents, []byte(footer))
}

func (w *streamWriter) Clean() error {
	return nil
}

package writer

import (
	"fmt"
	"io"
)

type streamWriter struct {
	dest io.Writer
}

// NewStreamWriter returns a new Writer that writes to the passed io.Writer.
// Files written will be prefixed with a comment containing the file name.
// nil can be passed to discard all writes.
func NewStreamWriter(dest io.Writer) Writer {
	if dest == nil {
		dest = io.Discard
	}

	return &streamWriter{dest: dest}
}

// WriteFile implements Writer.WriteFile
func (w *streamWriter) WriteFile(name string, contents []byte) error {
	header := fmt.Sprintf("// %s\n", fileName(name))
	footer := "\n\n"

	return multiWrite(w.dest, []byte(header), contents, []byte(footer))
}

// Clean implements Writer.Clean
func (w *streamWriter) Clean() error {
	return nil
}

package writer

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStreamWriter_noDest(t *testing.T) {
	w := NewStreamWriter(nil)
	assert.NotNil(t, w)

	assert.Equal(t, io.Discard, w.(*streamWriter).dest)
	assert.NotEqual(t, os.Stdout, w.(*streamWriter).dest)
}

func TestNewStreamWriter_dest(t *testing.T) {
	buf := &bytes.Buffer{}

	w := NewStreamWriter(buf)
	assert.NotNil(t, w)
	assert.Equal(t, buf, w.(*streamWriter).dest)
}

func TestStreamWriter_Clean(t *testing.T) {
	buf := &bytes.Buffer{}

	w := NewStreamWriter(buf)
	assert.NotNil(t, w)

	err := w.Clean()
	assert.NoError(t, err)
	assert.Empty(t, buf.Bytes())
}

func TestStreamWriter_WriteFile(t *testing.T) {
	buf := &bytes.Buffer{}

	w := NewStreamWriter(buf)
	assert.NotNil(t, w)

	err := w.WriteFile("foo", []byte("bar"))
	assert.NoError(t, err)

	expected := []byte("// zz_foo.go\nbar\n\n")
	assert.Equal(t, expected, buf.Bytes())
}

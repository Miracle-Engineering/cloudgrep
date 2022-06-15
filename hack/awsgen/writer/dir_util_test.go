package writer

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func testingDirWriter(t *testing.T) (Writer, string) {
	t.Helper()

	tmp := t.TempDir()
	w, err := NewDirWriter(tmp)

	require.NoError(t, err)
	require.NotNil(t, w)

	return w, tmp
}

func populateTestingDir(t *testing.T, dir string) []string {
	t.Helper()

	files := []struct {
		name     string
		contents string
		gen      bool
	}{
		{
			name:     "foo",
			contents: "bar",
		},
		{
			name:     "zz_foo",
			contents: "bar",
		},
		{
			name:     "foo.go",
			contents: "bar",
		},
		{
			name:     "zz_spam.go",
			contents: "ham",
			gen:      true,
		},
	}

	var nongen []string

	for _, f := range files {
		p := path.Join(dir, f.name)
		err := ioutil.WriteFile(p, []byte(f.contents), 0755)
		require.NoError(t, err)

		if !f.gen {
			nongen = append(nongen, f.name)
		}
	}

	return nongen
}

func dirFiles(t *testing.T, dir string) []string {
	t.Helper()
	var names []string

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	return names
}

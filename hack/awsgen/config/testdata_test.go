package config

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata
var testdata embed.FS

func prepTestdataDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	entries, err := testdata.ReadDir(".")
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		copyEntry(t, testdata, entry, ".", dir)
	}

	return path.Join(dir, "testdata")
}

func copyEntry(t *testing.T, srcFS embed.FS, srcEntry fs.DirEntry, srcDir, destDir string) {
	t.Helper()

	srcPath := path.Join(srcDir, srcEntry.Name())
	destPath := path.Join(destDir, srcEntry.Name())
	if srcEntry.IsDir() {
		t.Logf("copying dir %s to %s\n", srcPath, destPath)
		err := os.Mkdir(destPath, 0755)
		require.NoError(t, err)

		entries, err := srcFS.ReadDir(srcPath)
		require.NoError(t, err)

		for _, entry := range entries {
			copyEntry(t, srcFS, entry, srcPath, destPath)
		}
	} else {
		t.Logf("copying file %s to %s\n", srcPath, destPath)
		contents, err := srcFS.ReadFile(srcPath)
		require.NoError(t, err)

		err = ioutil.WriteFile(destPath, contents, 0644)
		require.NoError(t, err)
	}
}

package writer

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDirWriter_good(t *testing.T) {
	tmp := t.TempDir()

	w, err := NewDirWriter(tmp)
	assert.NoError(t, err)
	assert.NotNil(t, w)
}

func TestNewDirWriter_supportsSymlink(t *testing.T) {
	tmp := t.TempDir()

	real := path.Join(tmp, "real")
	sym := path.Join(tmp, "sym")

	err := os.Mkdir(real, 0755)
	require.NoError(t, err)

	err = os.Symlink(real, sym)
	require.NoError(t, err)

	w, err := NewDirWriter(sym)
	assert.NoError(t, err)
	assert.NotNil(t, w)
}

func TestNewDirWriter_notExists(t *testing.T) {
	tmp := t.TempDir()

	missing := path.Join(tmp, "missing")

	w, err := NewDirWriter(missing)
	assert.ErrorContains(t, err, "cannot read dir")
	assert.Nil(t, w)
}

func TestNewDirWriter_notDir(t *testing.T) {
	tmp := t.TempDir()

	filePath := path.Join(tmp, "regular")
	f, err := os.Create(filePath)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	w, err := NewDirWriter(filePath)
	assert.ErrorContains(t, err, "not a directory")
	assert.Nil(t, w)
}

func TestDirWriter_WriteFile_success(t *testing.T) {
	w, dir := testingDirWriter(t)

	name := "test"
	contents := []byte("foo\n")

	err := w.WriteFile(name, contents)
	assert.NoError(t, err)

	expectedName := "zz_" + name + ".go"
	expectedPath := path.Join(dir, expectedName)

	f, err := os.Open(expectedPath)
	assert.NoError(t, err)
	require.NotNil(t, f)

	defer f.Close()

	actualContents, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, contents, actualContents)
}

func TestDirWriter_WriteFile_overwrite(t *testing.T) {
	w, dir := testingDirWriter(t)

	name := "test"
	expectedName := "zz_" + name + ".go"
	expectedPath := path.Join(dir, expectedName)

	origContents := []byte("foo bar spam ham")
	err := ioutil.WriteFile(expectedPath, origContents, 0755)
	require.NoError(t, err)

	contents := []byte("foo\n")
	err = w.WriteFile(name, contents)
	assert.NoError(t, err)

	actualContents, err := ioutil.ReadFile(expectedPath)
	assert.NoError(t, err)
	assert.Equal(t, contents, actualContents)
}

func TestDirWriter_WriteFile_cannotOpen(t *testing.T) {
	w, dir := testingDirWriter(t)

	name := "test"
	expectedName := "zz_test.go"
	expectedPath := path.Join(dir, expectedName)
	err := os.Mkdir(expectedPath, 0755)
	require.NoError(t, err)

	err = w.WriteFile(name, []byte("foo"))
	assert.ErrorContains(t, err, "cannot open "+expectedPath+" for writing")
}

func TestDirWriter_Clean_noWrites(t *testing.T) {
	w, dir := testingDirWriter(t)
	nongen := populateTestingDir(t, dir)
	expected := nongen

	err := w.Clean()
	assert.NoError(t, err)

	files := dirFiles(t, dir)
	assert.ElementsMatch(t, expected, files)
}

func TestDirWriter_Clean_newWrite(t *testing.T) {
	w, dir := testingDirWriter(t)
	nongen := populateTestingDir(t, dir)
	expected := append(nongen, "zz_bar.go")

	err := w.WriteFile("bar", []byte("ham"))
	require.NoError(t, err)

	err = w.Clean()
	assert.NoError(t, err)

	files := dirFiles(t, dir)
	assert.ElementsMatch(t, expected, files)
}

func TestDirWriter_Clean_overwrite(t *testing.T) {
	w, dir := testingDirWriter(t)
	nongen := populateTestingDir(t, dir)
	expected := append(nongen, "zz_spam.go")

	err := w.WriteFile("spam", []byte("foo"))
	require.NoError(t, err)

	err = w.Clean()
	assert.NoError(t, err)

	files := dirFiles(t, dir)
	assert.ElementsMatch(t, expected, files)

	err = w.Clean()
	assert.NoError(t, err)
	files = dirFiles(t, dir)
	assert.ElementsMatch(t, expected, files)
}

func TestDirWriter_Clean_dir(t *testing.T) {
	w, dir := testingDirWriter(t)
	nongen := populateTestingDir(t, dir)

	name := "zz_foo.go"
	err := os.Mkdir(path.Join(dir, name), 0755)
	require.NoError(t, err)

	expected := append(nongen, name)

	err = w.Clean()
	assert.NoError(t, err)

	files := dirFiles(t, dir)
	assert.ElementsMatch(t, expected, files)
}

func TestDirWriter_Clean_cannotRemove(t *testing.T) {
	w, dir := testingDirWriter(t)
	populateTestingDir(t, dir)

	err := os.Chmod(dir, 0555)
	require.NoError(t, err)
	defer os.Chmod(dir, 0755)

	err = w.Clean()
	assert.ErrorContains(t, err, "cannot remove")
}

func TestDirWriter_Clean_cannotList(t *testing.T) {
	w, dir := testingDirWriter(t)
	populateTestingDir(t, dir)

	err := os.Chmod(dir, 0444)
	require.NoError(t, err)
	defer os.Chmod(dir, 0755)

	err = w.Clean()
	assert.ErrorContains(t, err, "cannot list contents")
}

package cmd

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	dir := t.TempDir()
	configPath := path.Join(dir, "config.yaml")
	writeConfig(t, configPath, "services: []")

	args := []string{
		"--config",
		configPath,
	}

	err := Do(args)
	assert.NoError(t, err)
}

func TestDo_noConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := path.Join(dir, "config.yaml")

	args := []string{
		"--config",
		configPath,
	}

	err := Do(args)
	assert.ErrorContains(t, err, "--config does not point to a valid config file")
}

func TestDo_configError(t *testing.T) {
	dir := t.TempDir()
	configPath := path.Join(dir, "config.yaml")
	svcPath := path.Join(dir, "svc.yaml")

	writeConfig(t, configPath, "services: [svc]")
	writeConfig(t, svcPath, "servicePackage: Foo")

	args := []string{
		"--config",
		configPath,
	}

	err := Do(args)
	assert.ErrorContains(t, err, "servicePackage not valid")
}

func TestDo_badOutDir(t *testing.T) {
	dir := t.TempDir()
	outDir := path.Join(dir, "out")
	configPath := path.Join(dir, "config.yaml")
	writeConfig(t, configPath, "services: []")

	args := []string{
		"--config",
		configPath,
		"--output-dir",
		outDir,
	}

	err := Do(args)
	assert.ErrorContains(t, err, "invalid --output-dir")
}

func TestDo_outDir(t *testing.T) {
	dir := t.TempDir()
	outDir := t.TempDir()
	configPath := path.Join(dir, "config.yaml")
	writeConfig(t, configPath, "services: []")

	args := []string{
		"--config",
		configPath,
		"--output-dir",
		outDir,
	}

	err := Do(args)
	assert.NoError(t, err)
}

func writeConfig(t *testing.T, path, content string) {
	t.Helper()

	err := ioutil.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

}

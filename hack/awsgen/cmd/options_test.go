package cmd

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	flags, opts := prepOptions()

	require.True(t, opts.Format, "expected --format=true")

	dir := t.TempDir()
	configPath := path.Join(dir, "foo.yaml")
	ioutil.WriteFile(configPath, []byte("services: []"), 0644)

	err := flags.Parse([]string{"--config", configPath, "--format=false"})
	assert.NoError(t, err)
	assert.False(t, opts.Format, "expected --format=false")

	err = opts.Validate()
	assert.NoError(t, err)
}

func TestOptions_Validate_noFile(t *testing.T) {
	flags, opts := prepOptions()

	require.True(t, opts.Format, "expected --format=true")

	dir := t.TempDir()
	configPath := path.Join(dir, "foo.yaml")

	err := flags.Parse([]string{"--config", configPath})
	assert.NoError(t, err)

	err = opts.Validate()
	assert.ErrorContains(t, err, "--config does not point to a valid config file")
}

func prepOptions() (*pflag.FlagSet, *Options) {
	flags := pflag.NewFlagSet("awsgen", pflag.ContinueOnError)
	opts := NewOptions(flags)

	return flags, opts
}

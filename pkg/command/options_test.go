package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	// Test default behavior
	opts, err := ParseOptions([]string{})
	assert.NoError(t, err)
	assert.Equal(t, false, opts.Debug)
	assert.Equal(t, "localhost", opts.HTTPHost)
	assert.Equal(t, uint(8080), opts.HTTPPort)

	// Test url prefix
	opts, err = ParseOptions([]string{"--prefix", "cloudgrep"})
	assert.NoError(t, err)
	assert.Equal(t, "cloudgrep/", opts.Prefix)

	opts, err = ParseOptions([]string{"--prefix", "cloudgrep/"})
	assert.NoError(t, err)
	assert.Equal(t, "cloudgrep/", opts.Prefix)

	opts, err = ParseOptions([]string{"--debug", "--bind", "0.0.0.0", "--listen", "8082", "--skip-open", "--version"})
	assert.NoError(t, err)
}

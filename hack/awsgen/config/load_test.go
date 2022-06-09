package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	dir := prepTestdataDir(t)

	rootPath := path.Join(dir, "config.yaml")
	cfg, err := Load(rootPath)
	assert.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Len(t, cfg.Services, 1)

	errs := cfg.Validate()
	assert.Len(t, errs, 0)
}

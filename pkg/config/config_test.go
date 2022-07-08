package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestInitDefault(t *testing.T) {
	config, err := GetDefault()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(config.Providers))
	assert.Equal(t, "aws", config.Providers[0].Cloud)

	assert.Equal(t, "sqlite", config.Datastore.Type)
	assert.Equal(t, "file::memory:?cache=shared", config.Datastore.DataSourceName)

	assert.Equal(t, "localhost", config.Web.Host)
	assert.Equal(t, 8080, config.Web.Port)
	assert.Equal(t, "/", config.Web.Prefix)
}

func TestReadFile(t *testing.T) {
	config, err := GetDefault()

	require.NoError(t, err)
	loadedConfig, err := ReadFile("test/dummy-config.yaml")
	require.NoError(t, err)
	require.Equal(t, "dummyhost", loadedConfig.Web.Host)
	config.Web.Host = loadedConfig.Web.Host
	require.Equal(t, config, loadedConfig)
}

func TestProfiles(t *testing.T) {
	//check that the no profile is set by default (would let AWS SDK decide based on env var)
	config, err := ReadFile("test/multi-region-config.yaml")
	require.NoError(t, err)
	err = config.Load()
	require.NoError(t, err)
	require.Equal(t, "", config.Providers[0].Profile)

	//check that each profile is declined in a dedicated provider
	config, _ = ReadFile("test/multi-region-config.yaml")
	config.Profiles = []string{"dev", "prod"}
	err = config.Load()
	require.NoError(t, err)
	require.Equal(t, 2, len(config.Providers))
	require.Equal(t, "dev", config.Providers[0].Profile)
	require.Equal(t, []string{"us-east-1", "global"}, config.Providers[0].Regions)
	require.Equal(t, "prod", config.Providers[1].Profile)
	require.Equal(t, []string{"us-east-1", "global"}, config.Providers[1].Regions)

	//check that profile can't be redefined if already set
	config, err = ReadFile("test/use-profile-config.yaml")
	require.NoError(t, err)
	config.Profiles = []string{"dev", "prod"}
	err = config.Load()
	require.ErrorContains(t, err, "using the option `--profiles` is not supported")

}

func TestRegions(t *testing.T) {
	// Manual regions trumps any written region.
	config, err := ReadFile("test/multi-region-config.yaml")
	require.NoError(t, err)
	require.Equal(t, 2, len(config.Providers[0].Regions))
	optionRegions := []string{"us-east-1"}
	config.Regions = optionRegions
	err = config.Load()
	require.NoError(t, err)
	require.Equal(t, optionRegions, config.Providers[0].Regions)

}

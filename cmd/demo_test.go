package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/run-x/cloudgrep/demo"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestDemoCommand(t *testing.T) {
	var actualConfig config.Config
	var actualLogger *zap.Logger

	dummyRunCmd := func(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
		actualConfig = cfg
		actualLogger = logger
		return nil
	}

	originalCmd := runCmd

	runCmd = dummyRunCmd
	demoConfig, _ := demo.GetDemoConfig()
	newPortConfig, _ := demo.GetDemoConfig()
	newHostConfig, _ := demo.GetDemoConfig()
	newHostConfig.Web.Host = "notlocalhost"
	newPortConfig.Web.Port = 8081
	var userConfig config.Config
	assert.NoError(t, yaml.Unmarshal(EmbedConfig, &userConfig))
	userConfig.Regions = []string{}

	testCases := []struct {
		name    string
		cfg     config.Config
		verbose bool
		args    []string
	}{
		{"AllGood", demoConfig, false, []string{"demo"}},
		{"AllGoodVerbose", demoConfig, true, []string{"demo", "-v"}},
		{"NewPort", newPortConfig, false, []string{"demo", "--port", "8081"}},
		{"NewPortShortHand", newPortConfig, false, []string{"demo", "-p", "8081"}},
		{"NewPortVerbose", newPortConfig, true, []string{"demo", "-v", "-p", "8081"}},
		{"NewHost", newHostConfig, false, []string{"demo", "--bind", "notlocalhost"}},
	}

	defer func() {
		runCmd = originalCmd
	}()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd := NewRootCmd(buf)
			rootCmd.SetArgs(tc.args)
			err := rootCmd.Execute()
			require.NotEqual(t, tc.cfg.Datastore.DataSourceName, actualConfig.Datastore.DataSourceName)
			tc.cfg.Datastore.DataSourceName = actualConfig.Datastore.DataSourceName
			require.NoError(t, err)
			require.Equal(t, tc.cfg, actualConfig)
			require.True(t, tc.verbose == actualLogger.Core().Enabled(zap.DebugLevel))
			require.Equal(t, 0, buf.Len())
		})
	}

}

package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"testing"
)

//go:embed test_config.yaml
var EmbedConfig []byte

func TestRootCommand(t *testing.T) {
	var actualConfig config.Config
	var actualLogger *zap.Logger

	dummyRunCmd := func(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
		actualConfig = cfg
		actualLogger = logger
		return nil
	}

	originalCmd := runCmd

	runCmd = dummyRunCmd
	defaultConfig, _ := config.GetDefault()
	newPortConfig, _ := config.GetDefault()
	newRegionsConfig, _ := config.GetDefault()
	newRegionsConfig.Regions = []string{"us-west-1", "us-west-2"}
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
		{"AllGood", defaultConfig, false, []string{}},
		{"AllGoodVerbose", defaultConfig, true, []string{"-v"}},
		{"Regions", newRegionsConfig, false, []string{"--regions", "us-west-1,us-west-2"}},
		{"RegionsShorthand", newRegionsConfig, false, []string{"-r", "us-west-1,us-west-2"}},
		{"NewPort", newPortConfig, false, []string{"--port", "8081"}},
		{"NewPortShortHand", newPortConfig, false, []string{"-p", "8081"}},
		{"NewPortVerbose", newPortConfig, true, []string{"-v", "-p", "8081"}},
		{"UserConfig", userConfig, false, []string{"-c", "./test_config.yaml"}},
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
			assert.NoError(t, err)
			assert.Equal(t, tc.cfg, actualConfig)
			assert.True(t, tc.verbose == actualLogger.Core().Enabled(zap.DebugLevel))
			assert.Equal(t, 0, buf.Len())
		})
	}

}

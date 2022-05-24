package cmd

import (
	"bytes"
	"context"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

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
	newPortConfig.Web.Port = 8081

	testCases := []struct {
		name    string
		cfg     config.Config
		verbose bool
		args    []string
	}{
		{"AllGood", defaultConfig, false, []string{}},
		{"AllGoodVerbose", defaultConfig, true, []string{"-v"}},
		{"NewPort", newPortConfig, false, []string{"--port", "8081"}},
		{"NewPortShortHand", newPortConfig, false, []string{"-p", "8081"}},
		{"NewPortVerbose", newPortConfig, true, []string{"-v", "-p", "8081"}},
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

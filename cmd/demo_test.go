package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDemoCommand(t *testing.T) {
	var actualConfig config.Config
	var actualLogger *zap.Logger

	originalCmd := runCmd
	runCmd = func(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
		actualConfig = cfg
		actualLogger = logger
		return nil
	}
	defer func() {
		runCmd = originalCmd
	}()

	testCases := []struct {
		name    string
		host    string
		port    int
		verbose bool
		args    []string
	}{
		{
			name:    "demo",
			host:    "localhost",
			port:    8080,
			verbose: false,
			args:    []string{"demo"},
		},
		{
			name:    "-v -p 8081",
			host:    "localhost",
			port:    8081,
			verbose: true,
			args:    []string{"demo", "-v", "-p", "8081"},
		},
		{
			name:    "--bind notlocalhost --port 8081",
			host:    "notlocalhost",
			port:    8081,
			verbose: false,
			args:    []string{"demo", "--bind", "notlocalhost", "--port", "8081"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd := NewRootCmd(buf)
			rootCmd.SetArgs(tc.args)
			err := rootCmd.Execute()
			require.NoError(t, err)
			//the datasource name should point to the demo data
			require.Contains(t, actualConfig.Datastore.DataSourceName, "cloudgrepdemodb")
			require.Equal(t, tc.host, actualConfig.Web.Host)
			require.Equal(t, tc.port, actualConfig.Web.Port)
			require.True(t, tc.verbose == actualLogger.Core().Enabled(zap.DebugLevel))
			require.Equal(t, 0, buf.Len())
		})
	}

}

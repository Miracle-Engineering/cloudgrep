package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"testing"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRootCommand(t *testing.T) {
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
		name     string
		host     string
		port     int
		regions  []string
		profiles []string
		roles    []string
		verbose  bool
		args     []string
	}{
		{
			name: "default",
			host: "localhost",
			port: 8080,
		},
		{
			name:    "-p 8081 -r us-west-1,us-west-2 -v",
			host:    "localhost",
			port:    8081,
			regions: []string{"us-west-1", "us-west-2"},
			verbose: true,
			args:    []string{"-p", "8081", "-v", "-r", "us-west-1,us-west-2"},
		},
		{
			name:     "--port 8081 --regions us-west-1,us-west-2 --profiles prod,dev --role-arns arn1,arn2",
			host:     "localhost",
			port:     8081,
			regions:  []string{"us-west-1", "us-west-2"},
			profiles: []string{"prod", "dev"},
			roles:    []string{"arn1", "arn2"},
			args:     []string{"--port", "8081", "--regions", "us-west-1,us-west-2", "--profiles", "prod,dev", "--role-arns", "arn1,arn2"},
		},
		{
			name:     "-r us-west-1 -r us-west-2 --profile dev --profile prod",
			host:     "localhost",
			port:     8080,
			regions:  []string{"us-west-1", "us-west-2"},
			profiles: []string{"dev", "prod"},
			args:     []string{"-r", "us-west-1", "-r", "us-west-2", "--profiles", "dev", "--profiles", "prod"},
		},
		{
			name: "-c ../pkg/config/test/custom-host-port.yaml",
			host: "helloworld",
			port: 8082,
			args: []string{"-c", "../pkg/config/test/custom-host-port.yaml"},
		},
		{
			name: "--port 8081 --config ../pkg/config/test/custom-host-port.yaml",
			host: "helloworld",
			port: 8081,
			args: []string{"--port", "8081", "--config", "../pkg/config/test/custom-host-port.yaml"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd := NewRootCmd(buf)
			rootCmd.SetArgs(tc.args)
			err := rootCmd.Execute()
			require.NoError(t, err)
			require.Equal(t, tc.regions, actualConfig.Regions)
			require.Equal(t, tc.profiles, actualConfig.Profiles)
			require.Equal(t, tc.roles, actualConfig.RoleArns)
			require.Equal(t, tc.host, actualConfig.Web.Host)
			require.Equal(t, tc.port, actualConfig.Web.Port)
			require.True(t, tc.verbose == actualLogger.Core().Enabled(zap.DebugLevel))
			require.Equal(t, 0, buf.Len())
		})
	}

}

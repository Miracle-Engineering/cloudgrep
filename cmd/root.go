package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/run-x/cloudgrep/pkg/cli"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger
var runCmd = cli.Run

type rootOptions struct {
	bind        string
	regions     string
	port        int
	prefix      string
	skipOpen    bool
	skipRefresh bool
	config      string
	verbose     bool
}

func (rO rootOptions) loadConfig() (config.Config, error) {
	var cfg config.Config
	var err error
	if rO.config != "" {
		cfg, err = config.LoadFromFile(rO.config)
	} else {
		cfg, err = config.GetDefault()
	}
	if err != nil {
		return cfg, err
	}
	if rO.bind != "" {
		cfg.Web.Host = rO.bind
	}
	if rO.regions != "" {
		cfg.Regions = strings.Split(rO.regions, ",")
	}
	if rO.port != 0 {
		cfg.Web.Port = rO.port
	}
	if rO.prefix != "" {
		cfg.Web.Prefix = ""
	}
	if rO.skipOpen {
		cfg.Web.SkipOpen = true
	}
	if rO.skipRefresh {
		cfg.Datastore.SkipRefresh = true
	}
	return cfg, nil
}

// NewRootCmd returns the base command when called without any subcommands
func NewRootCmd(out io.Writer) *cobra.Command {
	rO := rootOptions{}
	var rootCmd = &cobra.Command{
		Use:   "cloudgrep",
		Short: "A web-based utility to query and manage cloud resources",
		Long: `Cloudgrep is an app built by RunX to help devops manage the multitude of resources in
their cloud accounts.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := rO.loadConfig()
			if err != nil {
				return err
			}
			return runCmd(cmd.Context(), cfg, logger)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			if rO.verbose {
				logger, err = zap.NewDevelopment()
				util.EnableErrorStackTrace()
			} else {
				logger, err = zap.NewProduction()
			}
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&rO.verbose, "verbose", "v", false, "Log verbosity")

	flags := rootCmd.Flags()

	flags.StringVarP(&rO.config, "config", "c", "", "Config file (default is https://github.com/run-x/cloudgrep/blob/main/pkg/config/config.yaml)")
	flags.StringVar(&rO.bind, "bind", "", "Host to bind on")
	flags.StringVarP(&rO.regions, "regions", "r", "", "Comma separated list of regions to scan, or \"all\"")
	flags.IntVarP(&rO.port, "port", "p", 0, "Port to use")
	flags.StringVar(&rO.prefix, "prefix", "", "URL prefix to use")
	flags.BoolVar(&rO.skipOpen, "skip-open", false, "Skip running the open command to open default browser")
	flags.BoolVar(&rO.skipRefresh, "skip-refresh", false, "Skip running data refresh on start up")

	rootCmd.AddCommand(NewVersionCommand(out), NewDemoCommand(out))
	rootCmd.Commands()
	return rootCmd
}

func Execute() {
	err := NewRootCmd(os.Stdout).Execute()
	if err != nil {
		util.PrintStackTrace(err, os.Stderr)
		os.Exit(1)
	}
}

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/run-x/cloudgrep/pkg/cli"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var logger *zap.Logger
var runCmd = cli.Run

// NewRootCmd returns the base command when called without any subcommands
func NewRootCmd(out io.Writer) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "cloudgrep",
		Short: "A web-based utility to query and manage cloud resources",
		Long: `Cloudgrep is an app built by RunX to help devops manage the multitude of resources in
their cloud accounts.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg config.Config
			err := viper.Unmarshal(&cfg)
			if err != nil {
				return err
			}
			err = runCmd(cmd.Context(), cfg, logger)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.cloudgrep.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Log verbosity")

	defaultConfig, _ := config.GetDefault()
	flags := rootCmd.Flags()

	flags.String("bind", defaultConfig.Web.Host, "Host to bind on")
	_ = viper.BindPFlag("web.host", flags.Lookup("bind"))

	flags.StringP("regions", "r", "", "Comma separated list of regions to scan, or \"all\"")
	_ = viper.BindPFlag("regions", flags.Lookup("regions"))

	flags.IntP("port", "p", defaultConfig.Web.Port, "Port to use")
	_ = viper.BindPFlag("web.port", flags.Lookup("port"))

	flags.String("prefix", defaultConfig.Web.Prefix, "URL prefix to use")
	_ = viper.BindPFlag("web.prefix", flags.Lookup("prefix"))

	flags.Bool("skip-open", defaultConfig.Web.SkipOpen, "Skip running the open command to open default browser")
	_ = viper.BindPFlag("web.skipOpen", flags.Lookup("skip-open"))

	flags.Bool("skip-refresh", defaultConfig.Datastore.SkipRefresh, "Skip running data refresh on start up")
	_ = viper.BindPFlag("datastore.skipRefresh", flags.Lookup("skip-refresh"))

	rootCmd.AddCommand(NewVersionCommand(out))

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := NewRootCmd(os.Stdout).Execute()
	if err != nil {
		util.PrintStackTrace(err, os.Stderr)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cloudgrep" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cloudgrep")
	}
	var err error
	if verbose {
		logger, err = zap.NewDevelopment()
		util.EnableErrorStackTrace()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(config.EmbedConfig)); err != nil {
		panic(fmt.Errorf("Could not load default config"))
	}
	viper.SetEnvPrefix("cloudgrep")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	yamlData, err := yaml.Marshal(&cfg)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	logger.Sugar().Debugf("Using the following config \n%+v", string(yamlData))
}

// Package cmd
/*
Copyright © 2022 RunX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/run-x/cloudgrep/pkg/cli"
	"github.com/run-x/cloudgrep/pkg/config"
	"go.uber.org/zap"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var logger *zap.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudgrep",
	Short: "A web-based utility to query and manage cloud resources",
	Long: `Cloudgrep is an app built by RunX to help devops manage the multitude of resources in
their cloud accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.Config
		err := viper.Unmarshal(&cfg)
		if err != nil {
			panic(err)
		}
		logger.Sugar().Debugf("Using the following config %+v", cfg)
		err = cli.Run(cmd.Context(), cfg, logger)
		if err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.cloudgrep.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "log verbosity")

	defaultConfig, _ := config.GetDefault()

	rootCmd.Flags().String("bind", defaultConfig.Web.Host, "Host to bind on")
	_ = viper.BindPFlag("web.host", rootCmd.Flags().Lookup("bind"))

	rootCmd.Flags().IntP("port", "p", defaultConfig.Web.Port, "Port to use")
	_ = viper.BindPFlag("web.port", rootCmd.Flags().Lookup("port"))

	rootCmd.Flags().String("prefix", defaultConfig.Web.Prefix, "URL prefix to use")
	_ = viper.BindPFlag("web.prefix", rootCmd.Flags().Lookup("prefix"))

	rootCmd.Flags().Bool("skip-open", false, "Skip running the open command to open default browser")
	_ = viper.BindPFlag("web.skipOpen", rootCmd.Flags().Lookup("skip-open"))
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
}

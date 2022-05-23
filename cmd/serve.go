/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/run-x/cloudgrep/pkg/cli"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the data scan and turn on the web server",
	Long:  `Placeholder`,
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

func init() {
	rootCmd.AddCommand(serveCmd)
	defaultConfig, _ := config.GetDefault()

	serveCmd.Flags().String("bind", defaultConfig.Web.Host, "Host to bind on")
	_ = viper.BindPFlag("web.host", serveCmd.Flags().Lookup("bind"))

	serveCmd.Flags().IntP("port", "p", defaultConfig.Web.Port, "Port to use")
	_ = viper.BindPFlag("web.port", serveCmd.Flags().Lookup("port"))

	serveCmd.Flags().String("prefix", defaultConfig.Web.Prefix, "URL prefix to use")
	_ = viper.BindPFlag("web.prefix", serveCmd.Flags().Lookup("prefix"))

	serveCmd.Flags().Bool("skip-open", false, "Skip running the open command to open default browser")
	_ = viper.BindPFlag("web.skipOpen", serveCmd.Flags().Lookup("skip-open"))
}

package cmd

import (
	"context"
	"os"

	"github.com/run-x/cloudgrep/demo"
	"github.com/spf13/cobra"
)

type demoOptions struct {
	bind string
	port int
}

func (dO *demoOptions) run(ctx context.Context) error {
	cfg, err := demo.GetDemoConfig()
	if err != nil {
		return err
	}
	//clean up the temporary file
	defer os.Remove(cfg.Datastore.DataSourceName)
	if dO.bind != "" {
		cfg.Web.Host = dO.bind
	}
	if dO.port != 0 {
		cfg.Web.Port = dO.port
	}
	return runCmd(ctx, cfg, logger)
}

// NewDemoCommand returns the demo subcommand
func NewDemoCommand() *cobra.Command {
	var dO demoOptions
	var demoCmd = &cobra.Command{
		Use:   "demo",
		Short: "Run cloudgrep with a demo",
		Long: `The demo command runs demo with a pre-built config and database which shows some sample data for a medium-sized AWS account (a few hundred cloud resources).
This demo demonstrates how cloudgrep can help with:

- Viewing all your cloud resources for multiple regions in one browser.
- Searching your cloud resources using their tags to measure the progress of your IaC initiative.
- Verifying that your tag values are correct, quickly identifying the misconfigured values.
- Enforcing your tag policies by identifying the resources missing some tags.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return dO.run(cmd.Context())
		},
	}

	flags := demoCmd.Flags()
	flags.StringVar(&dO.bind, "bind", "", "Host to bind on")
	flags.IntVarP(&dO.port, "port", "p", 0, "Port to use")
	return demoCmd
}

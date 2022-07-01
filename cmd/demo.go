package cmd

import (
	"context"
	"github.com/run-x/cloudgrep/demo"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

type demoOptions struct {
	bind string
	port int
}

func (dO *demoOptions) run(ctx context.Context, out io.Writer) error {
	file, err := os.CreateTemp("", "cloudgrepdemodb")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	logger.Sugar().Infof("writing temporary file to store demo db: %v", file.Name())
	_, err = file.Write(demo.DemoDB)
	if err != nil {
		return err
	}
	cfg, err := demo.GetDemoConfig()
	if err != nil {
		return err
	}
	if dO.bind != "" {
		cfg.Web.Host = dO.bind
	}
	if dO.port != 0 {
		cfg.Web.Port = dO.port
	}
	cfg.Datastore.DataSourceName, err = filepath.Abs(file.Name())
	if err != nil {
		return err
	}
	return runCmd(ctx, cfg, logger)
}

// NewDemoCommand returns the demo subcommand
func NewDemoCommand(out io.Writer) *cobra.Command {
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
			return dO.run(cmd.Context(), out)
		},
	}

	flags := demoCmd.Flags()
	flags.StringVar(&dO.bind, "bind", "", "Host to bind on")
	flags.IntVarP(&dO.port, "port", "p", 0, "Port to use")
	return demoCmd
}

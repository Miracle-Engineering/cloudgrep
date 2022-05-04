package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/run-x/cloudgrep/pkg/api"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/options"
	"github.com/run-x/cloudgrep/pkg/provider"
	"github.com/run-x/cloudgrep/pkg/util"
)

func Run() error {
	ctx := context.Background()

	opts, err := options.ParseOptions(os.Args)
	if err != nil {
		return fmt.Errorf("failed to parse cli options: %w", err)
	}
	if opts.Version {
		fmt.Println(api.Version)
		os.Exit(0)
	}

	cfg, err := config.New(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Logging.IsDev() {
		util.StartProfiler()
	}

	//init the storage to contain cloud data
	datastore, err := datastore.NewDatastore(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to setup datastore: %w", err)
	}

	//start the providers to collect cloud data
	err = provider.Run(ctx, cfg, datastore)
	if err != nil {
		return fmt.Errorf("failed to start providers: %w", err)
	}

	api.StartServer(ctx, cfg, datastore)

	//TODO replace this URL with homepage when ready
	url := fmt.Sprintf("http://%v:%v/%v%v", cfg.Web.Host, cfg.Web.Port, cfg.Web.Prefix, "api/resources")
	fmt.Println("To view Cloudgrep UI, open ", url, "in browser")

	if !opts.SkipOpen {
		openPage(url)
	}
	handleSignals()
	return nil
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func openPage(url string) {
	_, err := exec.Command("which", "open").Output()
	if err != nil {
		return
	}

	err = exec.Command("open", url).Run()
	if err != nil {
		return
	}

}

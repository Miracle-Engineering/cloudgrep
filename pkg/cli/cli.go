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
	engine, err := provider.NewEngine(ctx, cfg, datastore)
	if err != nil {
		return fmt.Errorf("failed to start engine: %w", err)
	}
	if err = engine.Run(ctx); err != nil {
		stats, _ := datastore.Stats(ctx)
		if stats.ResourcesCount > 0 {
			//log the error but the api can still server with the datastore
			cfg.Logging.Logger.Sugar().Errorw("some error(s) when running the provider engine", "error", err)
		} else {
			// nothing to view - exit
			return fmt.Errorf("can't run the provider engine: %w", err)
		}
	}

	api.StartWebServer(ctx, cfg, datastore)

	url := fmt.Sprintf("http://%v:%v/%v%v", cfg.Web.Host, cfg.Web.Port, cfg.Web.Prefix, "api/fields")
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

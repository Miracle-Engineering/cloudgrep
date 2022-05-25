package cli

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/run-x/cloudgrep/pkg/api"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/provider"
	"github.com/run-x/cloudgrep/pkg/util"
)

func Run(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
	if logger.Core().Enabled(zap.DebugLevel) {
		util.StartProfiler()
	}

	//send amplitude event
	util.SendEvent(ctx, cfg, util.BaseEvent, nil)

	//init the storage to contain cloud data
	datastore, err := datastore.NewDatastore(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to setup datastore: %w", err)
	}

	//start the providers to collect cloud data
	engine, err := provider.NewEngine(ctx, cfg, logger, datastore)
	if err != nil {
		return fmt.Errorf("failed to start engine: %w", err)
	}
	if err = engine.Run(ctx); err != nil {
		stats, _ := datastore.Stats(ctx)
		if stats.ResourcesCount > 0 {
			//log the error but the api can still server with the datastore
			logger.Sugar().Errorw("some error(s) when running the provider engine", "error", err)
		} else {
			// nothing to view - exit
			return fmt.Errorf("can't run the provider engine: %w", err)
		}
	}

	api.StartWebServer(ctx, cfg, logger, datastore)

	url := fmt.Sprintf("http://%v:%v/%v", cfg.Web.Host, cfg.Web.Port, cfg.Web.Prefix)
	fmt.Println("To view Cloudgrep UI, open ", url, "in browser")

	if !cfg.Web.SkipOpen {
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

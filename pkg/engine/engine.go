package engine

import (
	"context"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider"
	"github.com/run-x/cloudgrep/pkg/sequencer"
	"log"

	"go.uber.org/zap"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
)

//Engine configures and starts the providers
type Engine struct {
	Providers []provider.Provider
	datastore.Datastore
	Logger    *zap.Logger
	Sequencer sequencer.Sequencer
}

//NewEngine Set up the providers, make sure configuration is valid
func NewEngine(ctx context.Context, cfg config.Config, logger *zap.Logger, datastore datastore.Datastore) (Engine, error) {
	e := Engine{}
	e.Datastore = datastore
	e.Logger = logger
	e.Sequencer = sequencer.AsyncSequencer{Logger: e.Logger}
	err := e.Datastore.WriteEvent(ctx, model.NewEngineEventStart())
	if err != nil {
		return e, err
	}
	for _, c := range cfg.Providers {
		// Manual regions trumps any written region.
		if len(cfg.Regions) > 0 {
			c.Regions = cfg.Regions
		}
		// create a providers
		err = datastore.WriteEvent(ctx, model.NewProviderEventStart(c.String()))
		if err != nil {
			log.Default().Println(err.Error())
			return e, err
		}
		providers, err := provider.NewProviders(ctx, c, logger)
		if err == nil {
			e.Providers = append(e.Providers, providers...)
		}
		err = datastore.WriteEvent(ctx, model.NewProviderEventEnd(c.String(), err))
		if err != nil {
			log.Default().Println(err.Error())
			return e, err
		}
	}
	return e, nil
}

//Run the providers: fetches data about cloud resources and save them to store
func (e *Engine) Run(ctx context.Context) error {
	err := e.Sequencer.Run(ctx, e, e.Providers)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	err = e.Datastore.WriteEvent(ctx, model.NewEngineEventLoaded())
	return err
}

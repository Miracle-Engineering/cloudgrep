package engine

import (
	"context"
	"github.com/run-x/cloudgrep/pkg/provider2"
	"github.com/run-x/cloudgrep/pkg/sequencer"

	"go.uber.org/zap"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
)

//Engine configures and starts the providers
type Engine struct {
	Providers []provider2.Provider
	datastore.Datastore
	Logger *zap.Logger
}

//Setup the providers, make sure configuration is valid
func NewEngine(ctx context.Context, cfg config.Config, logger *zap.Logger, datastore datastore.Datastore) (Engine, error) {
	e := Engine{}
	e.Datastore = datastore
	e.Logger = logger
	for _, config := range cfg.Providers {
		// create a provider
		providers, err := provider2.NewProviders(ctx, config, logger)
		if err != nil {
			return Engine{}, err
		}
		e.Providers = append(e.Providers, providers...)
	}
	return e, nil
}

//Run the providers: fetches data about cloud resources and save them to store
func (e *Engine) Run(ctx context.Context) error {
	asyncSequencer := sequencer.AsyncSequencer{Logger: e.Logger}
	return asyncSequencer.Run(ctx, e, e.Providers)
}

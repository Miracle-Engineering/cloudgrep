package engine

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider"
	"github.com/run-x/cloudgrep/pkg/sequencer"
	"go.uber.org/zap"
	"log"

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
	var errors *multierror.Error
	for _, c := range cfg.Providers {
		// Manual regions trumps any written region.
		if len(cfg.Regions) > 0 {
			c.Regions = cfg.Regions
		}
		// create a providers
		err := datastore.WriteEvent(ctx, model.NewProviderEventStart(c.String()))
		errors = multierror.Append(errors, err)
		if err != nil {
			log.Default().Println(err.Error())
		}
		providers, err := provider.NewProviders(ctx, c, logger)
		if err == nil {
			e.Providers = append(e.Providers, providers...)
		}
		errors = multierror.Append(errors, err)
		err = datastore.WriteEvent(ctx, model.NewProviderEventEnd(c.String(), err))
		if err != nil {
			log.Default().Println(err.Error())
		}
		errors = multierror.Append(errors, err)
	}
	return e, errors.ErrorOrNil()
}

//Run the providers: fetches data about cloud resources and save them to store
func (e *Engine) Run(ctx context.Context) error {
	var multipleErrors *multierror.Error
	err := e.Sequencer.Run(ctx, e, e.Providers)
	multipleErrors = multierror.Append(multipleErrors, err)
	if err != nil {
		log.Default().Println(err.Error())
	}
	err = e.Datastore.WriteEvent(ctx, model.NewEngineEventEnd(err))
	multipleErrors = multierror.Append(multipleErrors, err)
	return multipleErrors.ErrorOrNil()
}

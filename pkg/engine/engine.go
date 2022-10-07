package engine

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/juandiegopalomino/cloudgrep/pkg/config"
	"github.com/juandiegopalomino/cloudgrep/pkg/datastore"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/provider"
	"github.com/juandiegopalomino/cloudgrep/pkg/sequencer"
	"github.com/juandiegopalomino/cloudgrep/pkg/util/amplitude"
	"go.uber.org/zap"
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
	var errors error
	for _, c := range cfg.Providers {
		if err := datastore.WriteEvent(ctx, model.NewProviderEventStart(c.String())); err != nil {
			errors = multierror.Append(errors, err)
		}
		// create the providers - one per region
		providers, err := provider.NewProviders(ctx, c, logger)
		if err == nil {
			e.Providers = append(e.Providers, providers...)
			//send one amplitude event per AWS account
			for _, p := range providers {
				amplitude.SendEvent(logger, amplitude.EventCloudConnection, map[string]string{"CLOUD_ID": p.AccountId()})
			}
		} else {
			errors = multierror.Append(errors, err)
		}
		if err = datastore.WriteEvent(ctx, model.NewProviderEventEnd(c.String(), err)); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return e, errors
}

//Run the providers: fetches data about cloud resources and save them to store
func (e *Engine) Run(ctx context.Context) error {
	var errors error
	err := e.Sequencer.Run(ctx, e, e.Providers)
	if err != nil {
		errors = multierror.Append(errors, err)
	}
	if err = e.Datastore.WriteEvent(ctx, model.NewEngineEventEnd(err)); err != nil {
		errors = multierror.Append(errors, err)
	}
	return errors
}

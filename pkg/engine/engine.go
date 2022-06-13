package engine

import (
	"context"
	"github.com/hashicorp/go-multierror"
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

// NewEngine Set up the providers, make sure configuration is valid
func NewEngine(ctx context.Context, cfg config.Config, logger *zap.Logger, datastore datastore.Datastore) (Engine, error) {
	e := Engine{}
	e.Datastore = datastore
	e.Logger = logger
	e.Sequencer = sequencer.AsyncSequencer{Logger: e.Logger}
	for _, c := range cfg.Providers {
		// create a providers
		var errors *multierror.Error
		providers, err := provider.NewProviders(ctx, c, logger)
		errors = multierror.Append(errors, err)
		updateError := e.Datastore.UpdateProviderStatus(ctx, model.NewProviderStatus(c.Cloud, nil, err))
		errors = multierror.Append(errors, updateError)
		if errors.ErrorOrNil() != nil {
			return Engine{}, errors.ErrorOrNil()
		}
		e.Providers = append(e.Providers, providers...)
	}
	return e, nil
}

//Run the providers: fetches data about cloud resources and save them to store
func (e *Engine) Run(ctx context.Context) error {
	e.Datastore.CaptureEngineStart(ctx)
	err := e.Sequencer.Run(ctx, e, e.Providers)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	err = e.Datastore.CaptureEngineEnd(ctx)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

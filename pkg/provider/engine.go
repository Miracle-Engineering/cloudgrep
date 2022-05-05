package provider

import (
	"context"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
)

//Engine configures and starts the providers
type Engine struct {
	Providers []Provider
	datastore.Datastore
}

//Setup the providers, make sure configuration is valid
func NewEngine(ctx context.Context, cfg config.Config, datastore datastore.Datastore) (Engine, error) {
	e := Engine{}
	e.Datastore = datastore
	for _, config := range cfg.Providers {
		// create a provider
		provider, err := NewProvider(ctx, config, cfg.Logging.Logger)
		if err != nil {
			return Engine{}, err
		}
		e.Providers = append(e.Providers, provider)
	}
	return e, nil
}

//Run the providers: fetches data about cloud resources and save them to store
func (e *Engine) Run(ctx context.Context) error {

	//TODO use go routine to start the provider, review error handling to continue on error
	for _, provider := range e.Providers {
		// fetch the resources
		resources, err := fetchResources(ctx, provider)
		if err != nil {
			return err
		}
		// save to store
		err = e.WriteResources(ctx, resources)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchResources(ctx context.Context, provider Provider) ([]*model.Resource, error) {
	return provider.GetMapper().FetchResources(ctx, provider.Region())
}

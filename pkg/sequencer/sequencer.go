package sequencer

import (
	"context"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider2"
	"go.uber.org/zap"
	"sync"
)

type Sequencer interface {
	Run(ctx context.Context, ds datastore.Datastore, providers []provider2.Provider) error
}

type AsyncSequencer struct {
	Logger *zap.Logger
}

func (as AsyncSequencer) Run(ctx context.Context, ds datastore.Datastore, providers []provider2.Provider) error {
	resourceChan := make(chan model.Resource)
	doneChan := make(chan struct{})
	var wg sync.WaitGroup
	for _, provider := range providers {
		newFetchFuncs := provider.FetchFunctions()
		for resourceType, fetchFunc := range newFetchFuncs {
			wg.Add(1)
			go func(fetchFunc provider2.FetchFunc, provider provider2.Provider, resourceType string) {
				defer wg.Done()
				err := fetchFunc(ctx, resourceChan)
				if err != nil {
					// TODO: Log the error in like the future error table in db or somehow tell the user in the UI idk figure it out
					as.Logger.Sugar().Infof("Received an error when trying to handle resource %v in provider %v", resourceType, provider)
				}
			}(fetchFunc, provider, resourceType)
		}
	}

	go func() {
		defer close(doneChan)
		var resources []*model.Resource
	loop:
		for {
			select {
			case <-ctx.Done():
				return
			case resource, ok := <-resourceChan:
				if !ok {
					break loop
				}
				resources = append(resources, &resource)
			}
		}
		err := ds.WriteResources(ctx, resources)
		if err != nil {
			// TODO: Log the error in like the future error table in db or somehow tell the user in the UI idk figure it out
			as.Logger.Sugar().Infof("Received an error when trying to write resources to data store %v", err)
		}
	}()

	wg.Wait()
	close(resourceChan)
	<-doneChan
	return nil
}

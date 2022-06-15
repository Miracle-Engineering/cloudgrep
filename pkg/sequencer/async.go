package sequencer

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider"
	"go.uber.org/zap"
)

type AsyncSequencer struct {
	Logger *zap.Logger
}

func (as AsyncSequencer) Run(ctx context.Context, ds datastore.Datastore, providers []provider.Provider) error {
	resourceChan := make(chan model.Resource)
	doneChan := make(chan struct{})

	var errors *multierror.Error
	var errorLock sync.Mutex
	var wg sync.WaitGroup
	for _, p := range providers {
		newFetchFuncs := p.FetchFunctions()
		for resourceType, fetchFunc := range newFetchFuncs {
			wg.Add(1)
			go func(fetchFunc provider.FetchFunc, provider provider.Provider, resourceType string) {
				defer wg.Done()
				event := model.NewEvent(model.EventTypeResource, p.String(), resourceType)
				err := ds.WriteEvent(ctx, event)
				if err != nil {
					as.Logger.Sugar().Errorf("Received an error when trying to write resource event for resource  %v in provider %v: %v", resourceType, provider, err)
					errorLock.Lock()
					errors = multierror.Append(errors, err)
					errorLock.Unlock()
				}
				err = fetchFunc(ctx, resourceChan)
				event.UpdateError(err)
				err = ds.WriteEvent(ctx, event)
				if err != nil {
					as.Logger.Sugar().Errorf("Received an error when trying to write resource event for resource  %v in provider %v: %v", resourceType, provider, err)
					errorLock.Lock()
					errors = multierror.Append(errors, err)
					errorLock.Unlock()
				}
			}(fetchFunc, p, resourceType)
		}
	}

	var readError error
	go as.readResourceChan(ctx, doneChan, ds, resourceChan, &readError)

	wg.Wait()
	close(resourceChan)
	<-doneChan

	if readError != nil {
		errors = multierror.Append(errors, readError)
	}

	return errors.ErrorOrNil()
}

func (as AsyncSequencer) readResourceChan(ctx context.Context, doneCh chan<- struct{}, ds datastore.Datastore, resourceCh <-chan model.Resource, errOut *error) {
	defer close(doneCh)
	var resources []*model.Resource
loop:
	for {
		select {
		case <-ctx.Done():
			*errOut = ctx.Err()
			return
		case resource, ok := <-resourceCh:
			if !ok {
				break loop
			}
			resources = append(resources, &resource)
		}
	}
	err := ds.WriteResources(ctx, resources)
	if err != nil {
		// TODO: Log the error in like the future error table in db or somehow tell the user in the UI idk figure it out
		as.Logger.Sugar().Errorf("Received an error when trying to write resources to data store %v", err)
		*errOut = err
	}
}

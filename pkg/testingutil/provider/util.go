package provider

import (
	"context"
	"sync"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider/types"
	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func FetchResources[T types.Provider](ctx context.Context, t *testing.T, providers []T, name string) []model.Resource {
	t.Helper()

	var resources []model.Resource
	var resourceLock sync.Mutex
	var wg sync.WaitGroup

	worker := func(p types.Provider) {
		defer wg.Done()
		if p == nil {
			return
		}

		funcs := p.FetchFunctions()
		f, has := funcs[name]
		if !has {
			return
		}

		fetchFunc := testingutil.FetchAllFunc[model.Resource](f)

		funcResources, err := testingutil.FetchAll(ctx, t, fetchFunc)
		if err != nil {
			t.Errorf("failed to fetch %s on provider %s", name, p.String())
			return
		}

		for _, resource := range funcResources {
			resourceLock.Lock()
			resources = append(resources, resource)
			resourceLock.Unlock()
		}
	}

	for _, provider := range providers {
		wg.Add(1)
		go worker(provider)
	}

	wg.Wait()

	return resources
}

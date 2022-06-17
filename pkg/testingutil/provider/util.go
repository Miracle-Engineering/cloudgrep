package provider

import (
	"context"
	"sync"
	"sync/atomic"
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

	var foundCount int32

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

		atomic.AddInt32(&foundCount, 1)

		funcResources, err := testingutil.FetchAll(ctx, t, f)
		if err != nil {
			t.Errorf("failed to fetch %s on provider %s", name, p.String())
			return
		}

		funcResources = testingutil.ResourceFilterTagKeyValue(funcResources, "IntegrationTest", "true")

		if len(funcResources) > 0 {
			// Only count stats when we actually retrieve resources with populated tags
			stats.track(name)
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

	if foundCount == 0 {
		t.Fatalf("no providers found that define type %s", name)
	}

	return resources
}

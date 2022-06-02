package provider

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/assert"
)

func TestFakeProvider_String(t *testing.T) {
	fake := FakeProvider{}

	assert.Equal(t, "fake", fake.String())
	assert.Equal(t, "fake", fmt.Sprintf("%v", &fake))

	// Double check it implements fmt.Stringer properly
	var p provider.Provider = &fake
	assert.Equal(t, "fake", fmt.Sprintf("%v", p))
}

func TestFakeProvider_types(t *testing.T) {
	fake := FakeProvider{}

	configs := make(map[*FakeProviderResourceConfig]struct{})

	for key, val := range fake.types() {
		t.Run(key, func(t *testing.T) {
			val := val
			assert.NotEmpty(t, key)
			assert.NotNil(t, val)

			if _, has := configs[val]; has {
				t.Errorf("found duplicate pointer to config on type %s", key)
			}
			configs[val] = struct{}{}
		})
	}
}

func TestFakeProvider_FetchFunctions_idAffectsType(t *testing.T) {
	fake := FakeProvider{}
	funcs := fake.FetchFunctions()

	assert.Contains(t, funcs, "fake.Foo")

	fake.ID = "spam"
	funcs = fake.FetchFunctions()
	assert.Contains(t, funcs, "spam.Foo")
}

func TestFakeProvider_FetchFunctions_gutCheck(t *testing.T) {
	fake := FakeProvider{}

	for key, val := range fake.FetchFunctions() {
		t.Run(key, func(t *testing.T) {
			val := val
			assert.NotEmpty(t, key)
			assert.NotNil(t, val)
		})
	}
}

func TestFakeProvider_FetchFunctions_call(t *testing.T) {
	fake := FakeProvider{}
	fake.Foo.Count = 1
	fake.Bar.Count = 1

	funcs := fake.FetchFunctions()
	fetchFoo, has := funcs["fake.Foo"]
	if !has {
		t.Fatal("fake.Foo fetch func not found")
	}

	fetchAllFunc := testingutil.FetchAllFunc[model.Resource](fetchFoo)
	resources, err := testingutil.FetchAll(context.Background(), t, fetchAllFunc)

	assert.NoError(t, err)
	assert.Len(t, resources, 1)

	resource := resources[0]

	assert.Equal(t, "fake-fake.Foo-0", resource.Id)
	assert.Equal(t, 1, fake.Foo.RunCount())
	assert.Equal(t, 1, fake.TotalRuns())
}

func TestFakeProvider_FetchFunctions_callCanceled(t *testing.T) {
	fake := FakeProvider{}
	fake.Foo.Count = 100
	fake.Foo.DelayBefore = 1 * time.Second

	funcs := fake.FetchFunctions()
	fetchFoo := funcs["fake.Foo"]

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fetchAllFunc := testingutil.FetchAllFunc[model.Resource](fetchFoo)
	resources, err := testingutil.FetchAll(ctx, t, fetchAllFunc)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, resources)
}

func TestFakeProvider_FetchFunctions_delayBefore(t *testing.T) {
	t.Parallel()

	fake := FakeProvider{}
	fake.Foo.Count = 1
	fake.Foo.DelayBefore = 3 * time.Second

	funcs := fake.FetchFunctions()
	fetchFoo := funcs["fake.Foo"]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resourceCh := make(chan model.Resource)
	var atomicCount int32
	go func() {
		for range resourceCh {
			atomic.AddInt32(&atomicCount, 1)
		}
	}()

	go func() {
		defer close(resourceCh)
		fetchFoo(ctx, resourceCh)
	}()

	time.Sleep(2*time.Second + 500*time.Millisecond)

	count := atomic.LoadInt32(&atomicCount)
	assert.Zero(t, count)
}

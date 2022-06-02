package sequencer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/run-x/cloudgrep/pkg/provider2"
	"github.com/run-x/cloudgrep/pkg/testingutil/datastore"
	providerutil "github.com/run-x/cloudgrep/pkg/testingutil/provider"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func newAsync(t testing.TB) *AsyncSequencer {
	l := zaptest.NewLogger(t)

	s := &AsyncSequencer{
		Logger: l,
	}

	return s
}

func TestAsyncRun_good(t *testing.T) {
	t.Parallel()

	async := newAsync(t)
	ds := &datastore.Blackhole{}
	fakeProviders, providers := makeProviders()

	fakeProviders[0].Foo.DelayBefore = 2 * time.Second

	err := async.Run(context.Background(), ds, providers)
	assert.NoError(t, err)

	// Crude test to make sure that async.Run didn't exit before all fetch funcs complete
	assert.Equal(t, 2, fakeProviders[0].TotalRuns())

	assert.Equal(t, 3, ds.Count())
}

func TestAsyncRun_canceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	async := newAsync(t)
	ds := &datastore.Blackhole{}
	_, providers := makeProviders()

	err := async.Run(ctx, ds, providers)
	assert.ErrorIs(t, err, context.Canceled)

	assert.Equal(t, 0, ds.Count())
}

func TestAsyncRun_writeError(t *testing.T) {
	async := newAsync(t)
	ds := &datastore.Blackhole{}
	_, providers := makeProviders()

	expectedErr := errors.New("foo")
	ds.SetWriteError(expectedErr)

	err := async.Run(context.Background(), ds, providers)
	assert.ErrorIs(t, err, expectedErr)

	assert.Equal(t, 0, ds.Count())
}

func TestAsyncRun_fetchError(t *testing.T) {
	async := newAsync(t)
	ds := &datastore.Blackhole{}
	fakeProviders, providers := makeProviders()

	expectedErr := errors.New("foo")
	fakeProviders[0].Foo.ErrorAfterCount = expectedErr

	err := async.Run(context.Background(), ds, providers)
	assert.ErrorIs(t, err, expectedErr)

	assert.Equal(t, 3, ds.Count())
}

func makeProviders() ([]*providerutil.FakeProvider, []provider2.Provider) {
	fakeProviders := []*providerutil.FakeProvider{
		{
			ID: "spam",
			Foo: providerutil.FakeProviderResourceConfig{
				Count: 1,
			},
		},
		{
			ID: "ham",
			Bar: providerutil.FakeProviderResourceConfig{
				Count: 2,
			},
		},
	}

	providers := make([]provider2.Provider, len(fakeProviders))
	for idx, provider := range fakeProviders {
		providers[idx] = provider
	}

	return fakeProviders, providers
}

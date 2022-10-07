package engine

import (
	"context"
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/config"
	"github.com/juandiegopalomino/cloudgrep/pkg/datastore"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/provider"
	providerutil "github.com/juandiegopalomino/cloudgrep/pkg/testingutil/provider"
	"github.com/juandiegopalomino/cloudgrep/pkg/util/amplitude"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewEngine(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	t.Run("NoProviders", func(t *testing.T) {
		cfg, err := config.GetDefault()
		require.NoError(t, err)
		datastoreConfigs := config.Datastore{
			Type:           "sqlite",
			DataSourceName: "file::memory:",
		}
		cfg.Datastore = datastoreConfigs

		ds, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
		require.NoError(t, err)
		cfg.Providers = []config.Provider{}
		e, err := NewEngine(ctx, cfg, logger, ds)
		require.NoError(t, err)
		require.Equal(t, 0, len(e.Providers))
		require.Equal(t, logger, e.Logger)
	})

	t.Run("BadProvider", func(t *testing.T) {
		cfg, err := config.GetDefault()
		require.NoError(t, err)
		datastoreConfigs := config.Datastore{
			Type:           "sqlite",
			DataSourceName: "file::memory:",
		}
		cfg.Datastore = datastoreConfigs

		ds, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
		require.NoError(t, err)
		err = ds.WriteEvent(ctx, model.NewEngineEventStart())
		require.NoError(t, err)
		cfg.Providers = []config.Provider{{Cloud: "badCloud"}}
		engine, err := NewEngine(ctx, cfg, logger, ds)
		require.Error(t, err)
		engineEvent, err := engine.Datastore.EngineStatus(ctx)
		assert.Equal(t, model.EventStatusFetching, engineEvent.Status)
		assert.Equal(t, model.EventStatusFailed, engineEvent.ChildEvents[0].Status)
		require.NoError(t, err)
	})
}

type TestSequencer struct {
	Ran bool
}

func (ts *TestSequencer) Run(ctx context.Context, ds datastore.Datastore, providers []provider.Provider) error {
	ts.Ran = true
	return nil
}

func TestEngineRun(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	t.Run("RunEngine", func(t *testing.T) {
		cfg, err := config.GetDefault()
		require.NoError(t, err)
		datastoreConfigs := config.Datastore{
			Type:           "sqlite",
			DataSourceName: "file::memory:",
		}
		cfg.Datastore = datastoreConfigs

		//test events send to amplitude
		amplitudeClient := amplitude.UseMemoryClient()

		ds, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
		require.NoError(t, err)
		cfg.Providers = []config.Provider{
			{
				Cloud:   "fake",
				Regions: []string{"all}"},
			},
		}
		provider.RegisterExtraProviders("fake", fakeProviders())

		e, err := NewEngine(ctx, cfg, logger, ds)
		require.NoError(t, err)
		require.Equal(t, 2, len(e.Providers))
		require.Equal(t, logger, e.Logger)
		require.NoError(t, e.Run(ctx))

		//test resources created in the datastore
		resp, err := ds.GetResources(ctx, nil)
		assert.NoError(t, err)
		//the fake providers produce this number of resources
		assert.Equal(t, 3, len(resp.Resources))

		//check the datadtore events are correct
		engineStatus, err := ds.EngineStatus(ctx)
		assert.NoError(t, err)
		require.Equal(t, "success", engineStatus.Status)
		require.Equal(t, "", engineStatus.Error)

		//the engine sends one amplitude event per cloud id
		require.Equal(t, 2, amplitudeClient.Size())
		event, err := amplitudeClient.LastEvent()
		require.NoError(t, err)
		require.Equal(t, "CLOUD_CONNECTION", event["event_type"])
		require.Equal(t, "fake-provider-2", event["event_properties"].(map[string]string)["CLOUD_ID"])

	})
}

func fakeProviders() []provider.Provider {
	fakeProviders := []*providerutil.FakeProvider{
		{
			ID: "fake-provider-1",
			Foo: providerutil.FakeProviderResourceConfig{
				Count: 1,
			},
		},
		{
			ID: "fake-provider-2",
			Bar: providerutil.FakeProviderResourceConfig{
				Count: 2,
			},
		},
	}

	providers := make([]provider.Provider, len(fakeProviders))
	for idx, provider := range fakeProviders {
		providers[idx] = provider
	}

	return providers
}

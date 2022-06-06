package engine

import (
	"context"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/provider"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"testing"
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
		cfg.Providers = []config.Provider{{Cloud: "badCloud"}}
		_, err = NewEngine(ctx, cfg, logger, ds)
		require.Error(t, err)
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

		ds, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
		require.NoError(t, err)
		cfg.Providers = []config.Provider{}
		e, err := NewEngine(ctx, cfg, logger, ds)
		require.NoError(t, err)
		require.Equal(t, 0, len(e.Providers))
		require.Equal(t, logger, e.Logger)
		ts := &TestSequencer{Ran: false}
		e.Sequencer = ts
		require.NoError(t, e.Run(ctx))
		require.True(t, ts.Ran)
	})
}

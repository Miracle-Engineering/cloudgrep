package provider

import (
	"context"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func newTestEngine(t *testing.T) Engine {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	cfg := config.Config{
		Datastore: config.Datastore{
			Type:           "sqlite",
			DataSourceName: "file::memory:",
		},
	}
	datastore, err := datastore.NewSQLiteStore(ctx, cfg, logger)
	assert.NoError(t, err)

	engine, err := NewEngine(ctx, cfg, logger, datastore)
	assert.NoError(t, err)
	testProvider, err := NewTestProvider(ctx, config.Provider{}, logger)
	assert.NoError(t, err)
	engine.Providers = []Provider{testProvider}
	assert.NoError(t, err)
	return engine
}

func TestEngineRun(t *testing.T) {
	ctx := context.Background()
	//set some resources to return
	tr1 := TestResource{
		InstanceId:   "i-121",
		Architecture: nil,
		SomeTags:     []TestTag{},
	}
	tr2 := TestResource{
		InstanceId:   "i-122",
		Architecture: nil,
		SomeTags:     []TestTag{},
	}
	trMissingId := TestResource{
		Architecture: nil,
		SomeTags:     []TestTag{},
	}

	//run an engine that fetch 2 resources - no error
	engine := newTestEngine(t)
	err := engine.Run(
		context.WithValue(ctx, Return("FetchTestResources"), []TestResource{tr1, tr2}),
	)
	assert.NoError(t, err)
	engineStatus, err := engine.Datastore.GetEngineStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, model.EngineStatusSuccess, engineStatus.Status)
	//check that the resources were stored
	resources, err := engine.GetResources(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resources))

	//run an engine that fetch 2 resources - one with an error
	engine = newTestEngine(t)
	ctx = context.WithValue(ctx, Return("FetchTestResources"), []TestResource{tr1, trMissingId})
	err = engine.Run(ctx)
	assert.ErrorContains(t, err, "could not find id field")
	//check that 1 resource was stored
	resources, err = engine.GetResources(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}

func TestEngineRunResourceGetFailure(t *testing.T) {
	ctx := context.Background()
	//run an engine returns error
	engine := newTestEngine(t)
	err := engine.Run(
		context.WithValue(ctx, ReturnError("ReturnError"), "FetchTestError"),
	)
	assert.Error(t, err)
	engineStatus, err := engine.Datastore.GetEngineStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, model.EngineStatusFailed, engineStatus.Status)
}

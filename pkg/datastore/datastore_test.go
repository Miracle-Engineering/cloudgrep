package datastore

import (
	"context"
	_ "embed"
	"testing"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore/testdata"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestDatastore(t *testing.T) {

	datastoreConfigs := []config.Datastore{
		{
			Type: "memory",
		},
		{
			Type:           "sqlite",
			DataSourceName: "file::memory:?cache=shared",
		},
	}

	for _, datastoreConfig := range datastoreConfigs {
		t.Run(datastoreConfig.Type, func(t *testing.T) {
			ctx := context.Background()
			cfg := config.Config{
				Datastore: datastoreConfig,
				Logging: config.Logging{
					Logger: zaptest.NewLogger(t),
				},
			}
			dataStore, err := NewDatastore(ctx, cfg)
			assert.NoError(t, err)

			resources := testdata.GetResources(t)

			err = dataStore.WriteResources(ctx, resources)
			assert.NoError(t, err)

			var resourcesRead []*model.Resource
			resourcesRead, err = dataStore.GetResources(ctx, model.NoFilter{})
			assert.NoError(t, err)
			util.AssertEqualsResources(t, resources, resourcesRead)
		})
	}

}

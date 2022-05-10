package datastore

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore/testdata"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

var datastoreConfigs []config.Datastore

func newDatastores(t *testing.T, ctx context.Context) []Datastore {
	datastoreConfigs = []config.Datastore{
		{
			Type: "memory",
		},
		{
			Type:           "sqlite",
			DataSourceName: "file::memory:",
		},
	}
	var datastores []Datastore
	for _, datastoreConfig := range datastoreConfigs {
		cfg := config.Config{
			Datastore: datastoreConfig,
			Logging: config.Logging{
				Logger: zaptest.NewLogger(t),
				Mode:   "dev",
			},
		}
		dataStore, err := NewDatastore(ctx, cfg)
		assert.NoError(t, err)
		datastores = append(datastores, dataStore)
	}
	return datastores
}

func TestReadWrite(t *testing.T) {
	ctx := context.Background()
	for _, datastore := range newDatastores(t, ctx) {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			resources := testdata.GetResources(t)

			assert.NoError(t, datastore.WriteResources(ctx, resources))

			var resourcesRead []*model.Resource
			resourcesRead, err := datastore.GetResources(ctx, model.EmptyFilter())
			assert.NoError(t, err)
			util.AssertEqualsResources(t, resources, resourcesRead)
		})
	}
}
func TestFiltering(t *testing.T) {
	ctx := context.Background()
	for _, datastore := range newDatastores(t, ctx) {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			resources := testdata.GetResources(t)
			r1 := resources[0] //team:infra

			assert.NoError(t, datastore.WriteResources(ctx, resources))

			var resourcesRead []*model.Resource

			//only one resource has enabled=true
			filter := model.Filter{
				Tags: []model.Tag{{Key: "enabled", Value: "true"}},
			}
			resourcesRead, err := datastore.GetResources(ctx, filter)
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}
			//check 1 result returned
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead))
			util.AssertEqualsResourcePter(t, r1, resourcesRead[0])

			//check 2 tags filter: both resources have both tags - 2 results
			filter = model.Filter{
				Tags: []model.Tag{
					{Key: "vpc", Value: "vpc-123"},
					{Key: "eks:nodegroup", Value: "staging-default"},
				},
			}
			resourcesRead, err = datastore.GetResources(ctx, filter)
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead))
			util.AssertEqualsResources(t, resources, resourcesRead)

			//check 2 tags filter on same key - 2 results
			filter = model.Filter{
				Tags: []model.Tag{
					{Key: "team", Value: "infra"},
					{Key: "team", Value: "dev"},
				},
			}
			resourcesRead, err = datastore.GetResources(ctx, filter)
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead))
			util.AssertEqualsResources(t, resources, resourcesRead)

			//check 2 distinct tags - but no resource has both - 0 results
			filter = model.Filter{
				Tags: []model.Tag{
					{Key: "team", Value: "dev"},
					{Key: "env", Value: "prod"},
				},
			}
			resourcesRead, err = datastore.GetResources(ctx, filter)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(resourcesRead))

			//tag present - 2 results
			filter = model.Filter{
				Tags: []model.Tag{
					{Key: "team", Value: "*"},
				},
			}
			resourcesRead, err = datastore.GetResources(ctx, filter)
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead))
			util.AssertEqualsResources(t, resources, resourcesRead)

		})
	}
}

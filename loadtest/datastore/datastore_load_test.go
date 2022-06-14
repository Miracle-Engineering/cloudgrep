package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/datastore/testdata"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

var datastoreConfigs []config.Datastore

func newDatastores(t *testing.T, ctx context.Context) []datastore.Datastore {
	datastoreConfigs = []config.Datastore{
		{
			Type:           "sqlite",
			DataSourceName: "file::memory:",
		},
	}
	var datastores []datastore.Datastore
	for _, datastoreConfig := range datastoreConfigs {
		cfg := config.Config{
			Datastore: datastoreConfig,
		}
		dataStore, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
		assert.NoError(t, err)
		datastores = append(datastores, dataStore)
	}
	return datastores
}
func TestLoad(t *testing.T) {

	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	testCases := []struct {
		Resources int
		Tags      int
		BatchSize int
		Queries   int
		//below are generated
		Batches      int
		TagsPerBatch int
	}{
		//this is a smoke test to make sure the test works
		{
			Resources: 1,
			Tags:      1,
			BatchSize: 1,
			Queries:   1,
		},
		//this is the load test
		{
			Resources: 2000,
			//total includes the mandatory fields (region, type) and should be under 2000
			Tags:      1996,
			BatchSize: 100,
			Queries:   200,
		},
	}

	for _, tc := range testCases {
		tc.Batches = tc.Resources / tc.BatchSize
		tc.TagsPerBatch = tc.Tags / tc.Batches
		//total tags could be lower because of rounding
		tc.Tags = tc.TagsPerBatch * tc.Batches

		for _, ds := range newDatastores(t, ctx) {
			name := fmt.Sprintf("%T-%d-resources-%d-tags-%d-batch-%d-queries", ds, tc.Resources, tc.Tags, tc.BatchSize, tc.Queries)
			t.Run(name, func(t *testing.T) {

				test_resources := testdata.GetResources(t)
				assert.NotZero(t, len(test_resources))

				var wg sync.WaitGroup
				for i := 0; i < tc.Resources; i = i + tc.BatchSize {

					//for every batch generate some tag keys
					tags := make([]model.Tag, tc.TagsPerBatch)
					for t := range tags {
						tags[t] = model.Tag{
							Key:   fmt.Sprintf("batch-%d-tag-%d", i/tc.BatchSize, t),
							Value: uuid.New().String(),
						}
					}

					var newResources []*model.Resource
					for j := 0; j < tc.BatchSize; j = j + 1 {
						newResource := *test_resources[i%len(test_resources)]
						//set a unique id
						newResource.Id = fmt.Sprintf("i-%d", i*tc.BatchSize+j)
						//each tags instance should be different, otherwise gorm thinks it's the same
						newTags := make([]model.Tag, len(tags))
						copy(newTags, tags)
						newResource.Tags = newTags
						newResources = append(newResources, &newResource)

					}
					//each batch can be spawned into a routine to mirror async impl
					writeResources := func() {
						defer wg.Done()
						assert.NoError(t, ds.WriteResources(ctx, newResources))
					}
					wg.Add(1)
					go writeResources()

				}
				//wait insert to be done
				wg.Wait()

				if tc.Resources > datastore.DefaultLimit {
					//test limit is returned by default
					limitedResources, err := ds.GetResources(ctx, nil)
					assert.NoError(t, err)
					assert.True(t, len(limitedResources.Resources) < limitedResources.Count)
					assert.Equal(t, datastore.DefaultLimit, len(limitedResources.Resources))
				}

				//test all resources added
				allResources, err := ds.GetResources(ctx, []byte(fmt.Sprintf(`{"limit":%d}`, datastore.LimitMaxValue)))
				assert.NoError(t, err)
				if tc.Resources <= datastore.LimitMaxValue {
					assert.Equal(t, tc.Resources, len(allResources.Resources))
				}
				assert.Equal(t, tc.Resources, allResources.Count)
				//check the tags have been set
				assert.Equal(t, tc.TagsPerBatch, len(allResources.Resources[0].Tags))

				//test all fields
				fields, err := ds.GetFields(ctx)
				assert.NoError(t, err)
				assert.Equal(t, tc.Tags, len(fields.FindGroup("tags").Fields))

				//test some queries
				for i := 0; i < tc.Queries; i = i + 1 {
					//pick a random resource to find (to have at least 1 result)
					randResource := randResource(allResources.Resources)

					//build a rql query
					query := make(map[string]interface{})
					query["limit"] = tc.BatchSize
					//add some fields
					filter := make(map[string]interface{})
					for t, tag := range randResource.Tags {
						filter[tag.Key] = tag.Value
						if t == 0 {
							//add sort
							query["sort"] = []string{tag.Key}
						}
						//cap at 30 filter max
						if t == 30 {
							break
						}
					}
					query["filter"] = filter
					queryJson, err := json.Marshal(query)
					assert.NoError(t, err)
					result, err := ds.GetResources(ctx, queryJson)
					assert.NoError(t, err)
					//since all the tags are the same for a batch - the response should include them
					assert.Equal(t, tc.BatchSize, len(result.Resources))
				}
			})

		}
	}
}

func randResource(resources []*model.Resource) *model.Resource {
	return resources[rand.Intn(len(resources))]
}

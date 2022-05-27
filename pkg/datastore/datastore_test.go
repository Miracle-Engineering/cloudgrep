package datastore

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore/testdata"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

var datastoreConfigs []config.Datastore

const tagMaxKey = "service.k8s.aws/stack-XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDaFpLSjFbcXoEFfRsWxPLDnJObCsNVlgTeMaPEZQleQYhYzRyWJjPjzpfRFEgmotaFetHsbZRjxAwnwekrBEmfdzdcEkXBAkjQZLCtTMtTCoaNatyyiNKAReKJyiXJrscctNswYNsGRussVmaozFZBsbOJiFQGZsnwTKSmVoiGLOpbUOpEdKupdOMeRVjaRzL-----END"
const tagMaxValue = "ingress-nginx/ingress-nginx-controllerLDnJObCsNVlgTeMaPEZQleQYhYzRyWJjPjzpfRFEgmotaFetHsbZRjxAwnwekrBEEdKupdOMeRVjaRzL-----END"

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
		}
		dataStore, err := NewDatastore(ctx, cfg, zaptest.NewLogger(t))
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
			assert.NotZero(t, len(resources))

			//test write empty slice
			assert.NoError(t, datastore.WriteResources(ctx, []*model.Resource{}))

			//write the resources
			assert.NoError(t, datastore.WriteResources(ctx, resources))

			var resourcesRead []*model.Resource
			resourcesRead, err := datastore.GetResources(ctx, nil)
			assert.NoError(t, err)
			assert.Equal(t, len(resources), len(resourcesRead))
			model.AssertEqualsResources(t, resources, resourcesRead)

			//test getting a specific resource
			for _, r := range resources {
				resource, err := datastore.GetResource(ctx, r.Id)
				assert.NoError(t, err)
				model.AssertEqualsResourcePter(t, r, resource)
			}

		})
	}
}

func TestSearchByQuery(t *testing.T) {
	ctx := context.Background()
	for _, datastore := range newDatastores(t, ctx) {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			all_resources := testdata.GetResources(t)
			resourceInst1 := all_resources[0]  //i-123 team:infra, release tag, tag region:us-west-2
			resourceInst2 := all_resources[1]  //i-124 team:dev, no release tag
			resourceBucket := all_resources[2] //s3 bucket without tags

			assert.NoError(t, datastore.WriteResources(ctx, all_resources))

			var resourcesRead []*model.Resource

			//only one resource has enabled=true
			query := `{
  "filter":{
    "enabled": "true"
  }
}`

			resourcesRead, err := datastore.GetResources(ctx, []byte(query))
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}
			//check 1 result returned
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead))
			model.AssertEqualsResourcePter(t, resourceInst1, resourcesRead[0])

			//check 2 tags filter: both resources have both tags - 2 results
			query = `{
  "filter":{
    "vpc":"vpc-123",
    "eks:nodegroup":"staging-default"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead))
			model.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead)

			//check 2 tags filter on same key - 2 results
			query = `{
  "filter":{
    "$or":[
      {
        "team":"infra"
      },
      {
        "team":"dev"
      }
    ]
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead))
			model.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead)

			//check 2 distinct tags - but no resource has both - 0 result
			query = `{
  "filter":{
    "team":"dev",
    "env":"prod"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 0, len(resourcesRead))

			//tag present - 2 results
			query = `{
  "filter":{
	  "team": { "$neq": "" }
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead))
			model.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead)

			//test exclude - returns the resources without the tag release
			query = `{
  "filter":{
    "release": "[null]"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead))
			model.AssertEqualsResources(t, model.Resources{resourceInst2, resourceBucket}, resourcesRead)

			//test 2 exclusions - the s3 bucket is the only one without both tags
			query = `{
  "filter":{
    "release": "[null]",
    "debug:info": "[null]"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			model.AssertEqualsResources(t, model.Resources{resourceBucket}, resourcesRead)

			//mix include and exclude filters
			query = `{
  "filter":{
    "release":"[not null]",
    "vpc":"vpc-123"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead))
			model.AssertEqualsResourcePter(t, resourceInst1, resourcesRead[0])

			//test on max value
			query = fmt.Sprintf(`{"filter":{"%v":"%v"}}`, tagMaxKey, tagMaxValue)
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead))
			model.AssertEqualsResourcePter(t, resourceInst2, resourcesRead[0])

			//test on a tag called region - find the tag (ignore the core field)
			// we can probably revisit this in the future and include the group in the query field
			//ex: support "tags.region":"us-west-2" and "core.region":"us-west-2"
			query = `{
  "filter":{
    "region":"us-west-2"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead))
			model.AssertEqualsResourcePter(t, resourceInst1, resourcesRead[0])

		})
	}
}
func TestStats(t *testing.T) {
	ctx := context.Background()
	for _, datastore := range newDatastores(t, ctx) {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			resources := testdata.GetResources(t)
			assert.NoError(t, datastore.WriteResources(ctx, resources))

			stats, err := datastore.Stats(ctx)
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}
			//check stats
			assert.NoError(t, err)
			assert.Equal(t, model.Stats{ResourcesCount: 3}, stats)

		})
	}
}

func TestFields(t *testing.T) {
	ctx := context.Background()
	for _, datastore := range newDatastores(t, ctx) {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			resources := testdata.GetResources(t)
			assert.NoError(t, datastore.WriteResources(ctx, resources))

			fields, err := datastore.GetFields(ctx)
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}
			//check fields
			assert.NoError(t, err)
			assert.Equal(t, 11, len(fields))

			//test a few fields
			model.AssertEqualsField(t, model.Field{
				Group: "core",
				Name:  "region",
				Count: 3,
				Values: model.FieldValues{
					model.FieldValue{Value: "us-east-1", Count: 3},
				}}, *fields.Find("core", "region"))

			typeField := *fields.Find("core", "type")
			model.AssertEqualsField(t, model.Field{
				Group: "core",
				Name:  "type",
				Count: 3,
				Values: model.FieldValues{
					model.FieldValue{Value: "s3.Bucket", Count: 1},
					model.FieldValue{Value: "test.Instance", Count: 2},
				},
			}, typeField)

			//check that values are sorted by count desc
			assert.Equal(t, typeField.Values[0].Count, 2)
			assert.Equal(t, typeField.Values[1].Count, 1)

			model.AssertEqualsField(t, model.Field{
				Group: "tags",
				Name:  "team",
				Count: 2,
				Values: model.FieldValues{
					model.FieldValue{Value: "infra", Count: 1},
					model.FieldValue{Value: "dev", Count: 1},
				}}, *fields.Find("tags", "team"))

			//test long field
			model.AssertEqualsField(t, model.Field{
				Group: "tags",
				Name:  tagMaxKey,
				Count: 1,
				Values: model.FieldValues{
					model.FieldValue{Value: tagMaxValue, Count: 1},
				}}, *fields.Find("tags", tagMaxKey))

			//test the tag field called "region"
			model.AssertEqualsField(t, model.Field{
				Group: "tags",
				Name:  "region",
				Count: 1,
				Values: model.FieldValues{
					model.FieldValue{Value: "us-west-2", Count: 1},
				}}, *fields.Find("tags", "region"))

		})
	}
}

func TestLoad(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Only run on CI - a bit slow to run every time for developers")
	}

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

		for _, datastore := range newDatastores(t, ctx) {
			name := fmt.Sprintf("%T-%d-resources-%d-tags-%d-batch-%d-queries", datastore, tc.Resources, tc.Tags, tc.BatchSize, tc.Queries)
			t.Run(name, func(t *testing.T) {

				test_resources := testdata.GetResources(t)
				assert.NotZero(t, len(test_resources))

				//do not test result if not implemented
				_, err := datastore.GetFields(ctx)
				if err != nil && err.Error() == "not implemented" {
					return
				}

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
					assert.NoError(t, datastore.WriteResources(ctx, newResources))

				}
				if tc.Resources > defaultLimit {
					//test limit is returned by default
					limitedResources, err := datastore.GetResources(ctx, nil)
					assert.NoError(t, err)
					assert.Equal(t, defaultLimit, len(limitedResources))
				}

				//test all resources added
				allResources, err := datastore.GetResources(ctx, []byte(fmt.Sprintf(`{"limit":%d}`, limitMaxValue)))
				assert.NoError(t, err)
				if tc.Resources <= limitMaxValue {
					assert.Equal(t, tc.Resources, len(allResources))
				}
				//check the tags have been set
				assert.Equal(t, tc.TagsPerBatch, len(allResources[0].Tags))

				//test all fields
				fields, err := datastore.GetFields(ctx)
				assert.NoError(t, err)
				//the +2 is for region, type
				assert.Equal(t, tc.Tags+2, len(fields))

				//test some queries
				for i := 0; i < tc.Queries; i = i + 1 {
					//pick a random resource to find (to have at least 1 result)
					randResource := randResource(allResources)

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
					result, err := datastore.GetResources(ctx, queryJson)
					assert.NoError(t, err)
					//since all the tags are the same for a batch - the response should include them
					assert.Equal(t, tc.BatchSize, len(result))
				}
			})

		}
	}
}

func randResource(resources []*model.Resource) *model.Resource {
	return resources[rand.Intn(len(resources))]
}

package datastore

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore/testdata"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

//this is the max size allowed for an AWS tag
const tagMaxKey = "service.k8s.aws/stack-XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDaFpLSjFbcXoEFfRsWxPLDnJObCsNVlgTeMaPEZQleQYhYzRyWJjPjzpfRFEgmotaFetHsbZRjxAwnwekrBEmfdzdcEkXBAkjQZLCtTMtTCoaNatyyiNKAReKJyiXJrscctNswYNsGRussVmaozFZBsbOJiFQGZsnwTKSmVoiGLOpbUOpEdKupdOMeRVjaRzL-----END"
const tagMaxValue = "ingress-nginx/ingress-nginx-controllerLDnJObCsNVlgTeMaPEZQleQYhYzRyWJjPjzpfRFEgmotaFetHsbZRjxAwnwekrBEEdKupdOMeRVjaRzL-----END"

//only one resource has this tag
const tagUniqueResourceId = "i-123"
const tagUniqueKey = "unique-tag"
const tagUniqueValue = "unique-i-123"

func newDatastores(t *testing.T, ctx context.Context) ([]Datastore, []config.Config) {

	dbFilePath := path.Join(os.TempDir(), "cloudgrep-test.db")
	os.Remove(dbFilePath)

	datastoreConfigs := []config.Datastore{
		{
			Type:           "sqlite",
			DataSourceName: dbFilePath,
		},
	}
	var datastores []Datastore
	var configs []config.Config
	for _, datastoreConfig := range datastoreConfigs {
		cfg := config.Config{
			Datastore: datastoreConfig,
		}
		dataStore, err := NewDatastore(ctx, cfg, zaptest.NewLogger(t))
		assert.NoError(t, err)
		datastores = append(datastores, dataStore)
		configs = append(configs, cfg)
	}
	return datastores, configs
}

func TestReadWrite(t *testing.T) {
	ctx := context.Background()
	datastores, _ := newDatastores(t, ctx)
	for _, datastore := range datastores {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			resources := testdata.GetResources(t)
			assert.NotZero(t, len(resources))

			//test write empty slice
			assert.NoError(t, datastore.WriteResources(ctx, []*model.Resource{}))
			resourcesRead, err := datastore.GetResources(ctx, nil)
			assert.NoError(t, err)
			assert.Equal(t, model.ResourcesResponse{Count: 0, Resources: model.Resources{}}, resourcesRead)

			//write the resources
			assert.NoError(t, datastore.WriteResources(ctx, resources))

			resourcesRead, err = datastore.GetResources(ctx, nil)
			assert.NoError(t, err)
			assert.Equal(t, resourcesRead.Count, len(resourcesRead.Resources))
			assert.Equal(t, len(resources), len(resourcesRead.Resources))
			testingutil.AssertEqualsResources(t, resources, resourcesRead.Resources)

			//test getting a specific resource
			for _, r := range resources {
				resource, err := datastore.GetResource(ctx, r.Id)
				assert.NoError(t, err)
				testingutil.AssertEqualsResourcePter(t, r, resource)
			}

		})
	}
}

func TestSearchByQuery(t *testing.T) {
	ctx := context.Background()
	datastores, _ := newDatastores(t, ctx)
	for _, datastore := range datastores {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			all_resources := testdata.GetResources(t)
			resourceInst1 := all_resources[0]  //i-123 team:infra, release tag, tag region:us-west-2
			resourceInst2 := all_resources[1]  //i-124 team:dev, no release tag
			resourceBucket := all_resources[2] //s3 bucket without tags

			assert.NoError(t, datastore.WriteResources(ctx, all_resources))

			var resourcesRead model.ResourcesResponse

			//only one resource has enabled=true
			query := `{
  "filter":{
    "enabled": "true"
  }
}`

			resourcesRead, err := datastore.GetResources(ctx, []byte(query))
			//check 1 result returned
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead.Resources))
			assert.Equal(t, 1, resourcesRead.Count)
			testingutil.AssertEqualsResourcePter(t, resourceInst1, resourcesRead.Resources[0])

			//check 2 tags filter: both resources have both tags - 2 results
			query = `{
  "filter":{
    "vpc":"vpc-123",
    "eks:nodegroup":"staging-default"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead.Resources))
			assert.Equal(t, 2, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead.Resources)

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
			assert.Equal(t, 2, len(resourcesRead.Resources))
			assert.Equal(t, 2, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead.Resources)

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
  },
  "limit": 1
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead.Resources))
			assert.Equal(t, 2, resourcesRead.Count)

			//test 2 filters $or - both ec2 instances have these tags team and enabled
			//first $or returns 2 instances
			//second $or returns 1 instance --> result should be 1
			query = `{
  "filter":{
    "$or":[
      {
        "team":"infra"
      },
      {
        "team":"dev"
      }
    ],
        "$and": [
            { "$or": [
                { "enabled": "true" },
                { "enabled": "not-found" }
            ] }
        ]
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead.Resources))
			assert.Equal(t, 1, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{resourceInst1}, resourcesRead.Resources)

			//test 3 filter ors
			//1. "team":"(not null)" -> select both instances
			//2. "enabled": "(not null) -> select both instances
			//3. "id": "i-123" -> select 1 instance --> result should be 1
			query = `{
  "filter":{
    "$or":[
      {
        "team":"(not null)"
      }
    ],
    "$and":[
      {
        "$or":[
          {
            "enabled":"(not null)"
          }
        ]
      },
      {
        "$or":[
          {
            "id":"i-123"
          }
        ]
      }
    ]
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead.Resources))
			assert.Equal(t, 1, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{resourceInst1}, resourcesRead.Resources)

			//check 2 distinct tags - but no resource has both - 0 result
			query = `{
  "filter":{
    "team":"dev",
    "env":"prod"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 0, len(resourcesRead.Resources))
			assert.Equal(t, 0, resourcesRead.Count)

			//tag present - 2 results
			query = `{
  "filter":{
	  "team": { "$neq": "" }
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead.Resources))
			assert.Equal(t, 2, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead.Resources)

			//test exclude - returns the resources without the tag release
			query = `{
  "filter":{
    "release": "(null)"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resourcesRead.Resources))
			assert.Equal(t, 2, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{resourceInst2, resourceBucket}, resourcesRead.Resources)

			//test 2 exclusions - the s3 bucket is the only one without both tags
			query = `{
  "filter":{
    "release": "(null)",
    "debug:info": "(null)"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{resourceBucket}, resourcesRead.Resources)

			//mix include and exclude filters
			query = `{
  "filter":{
    "release":"(not null)",
    "vpc":"vpc-123"
  }
}`
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead.Resources))
			assert.Equal(t, 1, resourcesRead.Count)
			testingutil.AssertEqualsResourcePter(t, resourceInst1, resourcesRead.Resources[0])

			//test on max value
			query = fmt.Sprintf(`{"filter":{"%v":"%v"}}`, tagMaxKey, tagMaxValue)
			resourcesRead, err = datastore.GetResources(ctx, []byte(query))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead.Resources))
			assert.Equal(t, 1, resourcesRead.Count)
			testingutil.AssertEqualsResourcePter(t, resourceInst2, resourcesRead.Resources[0])

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
			assert.Equal(t, 1, len(resourcesRead.Resources))
			assert.Equal(t, 1, resourcesRead.Count)
			testingutil.AssertEqualsResourcePter(t, resourceInst1, resourcesRead.Resources[0])

		})
	}
}
func TestStats(t *testing.T) {
	ctx := context.Background()
	datastores, _ := newDatastores(t, ctx)
	for _, datastore := range datastores {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			resources := testdata.GetResources(t)
			assert.NoError(t, datastore.WriteResources(ctx, resources))

			stats, err := datastore.Stats(ctx)
			//check stats
			assert.NoError(t, err)
			assert.Equal(t, model.Stats{ResourcesCount: 3}, stats)

		})
	}
}

func TestFields(t *testing.T) {
	ctx := context.Background()
	datastores, _ := newDatastores(t, ctx)
	for _, datastore := range datastores {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			resources := testdata.GetResources(t)
			assert.NoError(t, datastore.WriteResources(ctx, resources))

			fields, err := datastore.GetFields(ctx)
			//check fields
			assert.NoError(t, err)
			//check number of groups
			assert.Equal(t, 2, len(fields))
			//check fields by group
			assert.Equal(t, 2, len(fields.FindGroup("core").Fields))
			assert.Equal(t, 10, len(fields.FindGroup("tags").Fields))

			//test a few fields
			testingutil.AssertEqualsField(t, model.Field{
				Name:  "region",
				Count: 3,
				Values: model.FieldValues{
					model.FieldValue{Value: "us-east-1", Count: 3},
				}}, *fields.FindField("core", "region"))

			typeField := *fields.FindField("core", "type")
			testingutil.AssertEqualsField(t, model.Field{
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

			testingutil.AssertEqualsField(t, model.Field{
				Name:  "team",
				Count: 2,
				Values: model.FieldValues{
					model.FieldValue{Value: "infra", Count: 1},
					model.FieldValue{Value: "dev", Count: 1},
					model.FieldValue{Value: "(null)", Count: 1},
				}}, *fields.FindField("tags", "team"))

			//test long field
			testingutil.AssertEqualsField(t, model.Field{
				Name:  tagMaxKey,
				Count: 1,
				Values: model.FieldValues{
					model.FieldValue{Value: tagMaxValue, Count: 1},
					model.FieldValue{Value: "(null)", Count: 2},
				}}, *fields.FindField("tags", tagMaxKey))

			//test the tag field called "region"
			testingutil.AssertEqualsField(t, model.Field{
				Name:  "region",
				Count: 1,
				Values: model.FieldValues{
					model.FieldValue{Value: "us-west-2", Count: 1},
					model.FieldValue{Value: "(null)", Count: 2},
				}}, *fields.FindField("tags", "region"))

		})
	}
}

//test that the resources can be updated: update their properties, tags
func TestUpdateResources(t *testing.T) {
	ctx := context.Background()
	datastores, _ := newDatastores(t, ctx)
	for _, ds := range datastores {
		name := fmt.Sprintf("%T", ds)
		t.Run(name, func(t *testing.T) {

			startTime := time.Now()

			//write and read the resources
			resources := testdata.GetResources(t)
			require.NoError(t, ds.WriteResources(ctx, resources))
			r1, err := ds.GetResource(ctx, resources[0].Id)
			require.NoError(t, err)

			//test the UpdatedAt field has been set
			require.NotNil(t, r1.UpdatedAt)
			require.Greater(t, r1.UpdatedAt, startTime)
			lastUpdatedAt := r1.UpdatedAt

			//update a resource - test the updated value and the timestamp was updated
			r1.Region = "us-west-2"
			require.NoError(t, ds.WriteResources(ctx, model.Resources{r1}))
			r1, err = ds.GetResource(ctx, resources[0].Id)
			require.NoError(t, err)
			require.Equal(t, r1.Region, "us-west-2")
			require.Greater(t, r1.UpdatedAt, lastUpdatedAt)

			//test updating some tags
			r1UniqueTag := model.Tag{Key: tagUniqueKey, Value: tagUniqueValue}
			deletedTag := r1UniqueTag
			//before deleting, check that the query on that tag returns the resource
			testQuery(t, ctx, ds, deletedTag.Key, deletedTag.Value, r1)
			//delete and add a new tag
			tags := r1.Tags.Delete(deletedTag.Key)
			newTag := model.Tag{Key: "brand-new", Value: "shinny"}
			tags = append(tags, newTag)
			r1.Tags = tags
			require.NoError(t, ds.WriteResources(ctx, model.Resources{r1}))
			r1, err = ds.GetResource(ctx, resources[0].Id)
			require.NoError(t, err)
			//test new tag is found
			testingutil.AssertEqualsTag(t, &newTag, r1.Tags.Find("brand-new"))
			testQuery(t, ctx, ds, newTag.Key, newTag.Value, r1)
			//test old tag is deleted
			testingutil.AssertEqualsTag(t, nil, r1.Tags.Find(deletedTag.Key))
			//searching on deleting tag returns nothing
			testQueryNoResult(t, ctx, ds, deletedTag.Key, deletedTag.Value)

		})
	}
}

//test that the DB can be reloaded at startup
func TestReloadDB(t *testing.T) {
	ctx := context.Background()
	datastores, configs := newDatastores(t, ctx)
	for i, ds := range datastores {
		cfg := configs[i]
		name := fmt.Sprintf("%T", ds)
		t.Run(name, func(t *testing.T) {

			//simulate a 1st run that would write resources to the datastore
			resources := testdata.GetResources(t)
			require.NoError(t, ds.WriteResources(ctx, resources))
			resourcesRead, err := ds.GetResources(ctx, nil)
			require.NoError(t, err)
			assert.Equal(t, len(resourcesRead.Resources), resourcesRead.Count)
			assert.NotZero(t, len(resourcesRead.Resources))
			r1, _ := ds.GetResource(ctx, tagUniqueResourceId)
			assert.NotNil(t, r1)
			//test a query
			testQuery(t, ctx, ds, tagUniqueKey, tagUniqueValue, r1)

			//simulate a 2nd run that would reload the resources (no write done)
			dsNew, err := NewDatastore(ctx, cfg, zaptest.NewLogger(t))
			require.NoError(t, err)
			resourcesReadNew, err := dsNew.GetResources(ctx, nil)
			require.NoError(t, err)
			require.Equal(t, len(resourcesReadNew.Resources), resourcesRead.Count)
			//the new datastore contains the same data that was previsouly stored
			testingutil.AssertEqualsResources(t, resourcesRead.Resources, resourcesReadNew.Resources)
			//test the same query - test index were loaded
			testQuery(t, ctx, dsNew, tagUniqueKey, tagUniqueValue, r1)

		})
	}
}

//test that the resources that no longer exist are purged
func TestPurgeResources(t *testing.T) {
	ctx := context.Background()
	datastores, _ := newDatastores(t, ctx)
	for _, ds := range datastores {
		name := fmt.Sprintf("%T", ds)
		t.Run(name, func(t *testing.T) {

			//1nd run: write 3 resources
			resources := testdata.GetResources(t)[:3]
			require.NoError(t, ds.WriteEvent(ctx, model.NewEngineEventStart()))
			require.NoError(t, ds.WriteResources(ctx, resources))
			require.NoError(t, ds.WriteEvent(ctx, model.NewEngineEventLoaded()))
			r1, err := ds.GetResource(ctx, resources[0].Id)
			require.NoError(t, err)
			r2, err := ds.GetResource(ctx, resources[1].Id)
			require.NoError(t, err)
			r3, err := ds.GetResource(ctx, resources[2].Id)
			require.NoError(t, err)
			testQuery(t, ctx, ds, tagUniqueKey, tagUniqueValue, r1)

			//2nd run: one resource is removed
			require.NoError(t, ds.WriteEvent(ctx, model.NewEngineEventStart()))
			require.NoError(t, ds.WriteResources(ctx, model.Resources{r2, r3}.Clean()))
			require.NoError(t, ds.WriteEvent(ctx, model.NewEngineEventLoaded()))
			resourcesRead, err := ds.GetResources(ctx, nil)
			require.NoError(t, err)
			require.Equal(t, 2, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{r2, r3}, resourcesRead.Resources)
			//the query doesn't return the deleted resource
			testQueryNoResult(t, ctx, ds, "id", r1.Id)
			testQueryUnrecognizedKey(t, ctx, ds, tagUniqueKey, tagUniqueValue)

			//3rd run: an error happened - there is a built-in protection to not delete all resources
			require.NoError(t, ds.WriteEvent(ctx, model.NewEngineEventStart()))
			require.NoError(t, ds.WriteEvent(ctx, model.NewEngineEventLoaded()))
			resourcesRead, err = ds.GetResources(ctx, nil)
			require.NoError(t, err)
			require.Equal(t, 2, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{r2, r3}, resourcesRead.Resources)
			testQuery(t, ctx, ds, "id", r2.Id, r2)

			//4th run: add back the resource previously deleted
			require.NoError(t, ds.WriteEvent(ctx, model.NewEngineEventStart()))
			require.NoError(t, ds.WriteResources(ctx, model.Resources{r2, r1, r3}.Clean()))
			require.NoError(t, ds.WriteEvent(ctx, model.NewEngineEventLoaded()))
			resourcesRead, err = ds.GetResources(ctx, nil)
			require.NoError(t, err)
			require.Equal(t, 3, resourcesRead.Count)
			testingutil.AssertEqualsResources(t, model.Resources{r1, r2, r3}, resourcesRead.Resources)
			testQuery(t, ctx, ds, tagUniqueKey, tagUniqueValue, r1)

		})
	}
}

func TestEngineStatus(t *testing.T) {
	type args struct {
		events model.Events
	}
	tests := []struct {
		name     string
		args     args
		expected model.Event
	}{
		{
			name: "success-status",
			args: args{
				events: model.Events{
					model.NewEngineEventStart(),
					model.NewProviderEventStart("mp1"),
					model.NewProviderEventStart("mp2"),
					model.NewProviderEventEnd("mp1", nil),
					model.NewProviderEventEnd("mp2", nil),
					model.NewEngineEventLoaded(),
					model.NewResourceEventStart("mp1", "mr1"),
					model.NewResourceEventStart("mp1", "mr2"),
					model.NewResourceEventStart("mp1", "mr3"),
					model.NewResourceEventEnd("mp1", "mr1", nil),
					model.NewResourceEventEnd("mp1", "mr2", nil),
					model.NewResourceEventEnd("mp1", "mr3", nil),
					model.NewResourceEventStart("mp2", "mr1"),
					model.NewResourceEventStart("mp2", "mr2"),
					model.NewResourceEventStart("mp2", "mr3"),
					model.NewResourceEventEnd("mp2", "mr1", nil),
					model.NewResourceEventEnd("mp2", "mr2", nil),
					model.NewResourceEventEnd("mp2", "mr3", nil),
				},
			},
			expected: model.Event{
				Type:   model.EventTypeEngine,
				Status: model.EventStatusSuccess,
				ChildEvents: model.Events{
					model.Event{
						Type:         model.EventTypeProvider,
						Status:       model.EventStatusSuccess,
						ProviderName: "mp1",
						ChildEvents: model.Events{
							model.NewResourceEventEnd("mp1", "mr1", nil),
							model.NewResourceEventEnd("mp1", "mr2", nil),
							model.NewResourceEventEnd("mp1", "mr3", nil),
						},
					},
					model.Event{
						Type:         model.EventTypeProvider,
						Status:       model.EventStatusSuccess,
						ProviderName: "mp2",
						ChildEvents: model.Events{
							model.NewResourceEventEnd("mp2", "mr1", nil),
							model.NewResourceEventEnd("mp2", "mr2", nil),
							model.NewResourceEventEnd("mp2", "mr3", nil),
						},
					},
				},
			},
		},
		{
			name: "failed-provider-status",
			args: args{
				events: model.Events{
					model.NewEngineEventStart(),
					model.NewProviderEventStart("mp1"),
					model.NewProviderEventStart("mp2"),
					model.NewProviderEventEnd("mp1", nil),
					model.NewProviderEventEnd("mp2", errors.New("mp2-error")),
					model.NewEngineEventLoaded(),
					model.NewResourceEventStart("mp1", "mr1"),
					model.NewResourceEventStart("mp1", "mr2"),
					model.NewResourceEventStart("mp1", "mr3"),
					model.NewResourceEventEnd("mp1", "mr1", nil),
					model.NewResourceEventEnd("mp1", "mr2", nil),
					model.NewResourceEventEnd("mp1", "mr3", nil),
				},
			},
			expected: model.Event{
				Type:   model.EventTypeEngine,
				Status: model.EventStatusFailed,
				Error:  "mp2-error",
				ChildEvents: model.Events{
					model.Event{
						Type:         model.EventTypeProvider,
						Status:       model.EventStatusSuccess,
						ProviderName: "mp1",
						ChildEvents: model.Events{
							model.NewResourceEventEnd("mp1", "mr1", nil),
							model.NewResourceEventEnd("mp1", "mr2", nil),
							model.NewResourceEventEnd("mp1", "mr3", nil),
						},
					},
					model.NewProviderEventEnd("mp2", errors.New("mp2-error")),
				},
			},
		},
		{
			name: "failed-provider-failed-resource-status",
			args: args{
				events: model.Events{
					model.NewEngineEventStart(),
					model.NewProviderEventStart("mp1"),
					model.NewProviderEventStart("mp2"),
					model.NewProviderEventEnd("mp1", nil),
					model.NewProviderEventEnd("mp2", errors.New("mp2-error")),
					model.NewEngineEventLoaded(),
					model.NewResourceEventStart("mp1", "mr1"),
					model.NewResourceEventStart("mp1", "mr2"),
					model.NewResourceEventStart("mp1", "mr3"),
					model.NewResourceEventEnd("mp1", "mr1", nil),
					model.NewResourceEventEnd("mp1", "mr2", nil),
					model.NewResourceEventEnd("mp1", "mr3", errors.New("mp1-mr3-error")),
				},
			},
			expected: model.Event{
				Type:   model.EventTypeEngine,
				Status: model.EventStatusFailed,
				Error:  "mp1-mr3-error\nmp2-error",
				ChildEvents: model.Events{
					model.Event{
						Type:         model.EventTypeProvider,
						Status:       model.EventStatusFailed,
						ProviderName: "mp1",
						Error:        "mp1-mr3-error",
						ChildEvents: model.Events{
							model.NewResourceEventEnd("mp1", "mr1", nil),
							model.NewResourceEventEnd("mp1", "mr2", nil),
							model.NewResourceEventEnd("mp1", "mr3", errors.New("mp1-mr3-error")),
						},
					},
					model.NewProviderEventEnd("mp2", errors.New("mp2-error")),
				},
			},
		},
		{
			name: "fetching-status",
			args: args{
				events: model.Events{
					model.NewEngineEventStart(),
					model.NewProviderEventStart("mp1"),
					model.NewProviderEventStart("mp2"),
					model.NewProviderEventEnd("mp1", nil),
					model.NewProviderEventEnd("mp2", nil),
					model.NewEngineEventLoaded(),
					model.NewResourceEventStart("mp1", "mr1"),
					model.NewResourceEventStart("mp1", "mr2"),
					model.NewResourceEventStart("mp1", "mr3"),
					model.NewResourceEventStart("mp1", "mr4"),
					model.NewResourceEventEnd("mp1", "mr1", nil),
					model.NewResourceEventEnd("mp1", "mr2", nil),
					model.NewResourceEventEnd("mp1", "mr3", nil),
					model.NewResourceEventStart("mp2", "mr1"),
					model.NewResourceEventStart("mp2", "mr2"),
					model.NewResourceEventStart("mp2", "mr3"),
					model.NewResourceEventEnd("mp2", "mr1", nil),
					model.NewResourceEventEnd("mp2", "mr2", nil),
					model.NewResourceEventEnd("mp2", "mr3", nil),
				},
			},
			expected: model.Event{
				Type:   model.EventTypeEngine,
				Status: model.EventStatusFetching,
				ChildEvents: model.Events{
					model.Event{
						Type:         model.EventTypeProvider,
						Status:       model.EventStatusFetching,
						ProviderName: "mp1",
						ChildEvents: model.Events{
							model.NewResourceEventEnd("mp1", "mr1", nil),
							model.NewResourceEventEnd("mp1", "mr2", nil),
							model.NewResourceEventEnd("mp1", "mr3", nil),
							model.NewResourceEventStart("mp1", "mr4"),
						},
					},
					model.Event{
						Type:         model.EventTypeProvider,
						Status:       model.EventStatusSuccess,
						ProviderName: "mp2",
						ChildEvents: model.Events{
							model.NewResourceEventEnd("mp2", "mr1", nil),
							model.NewResourceEventEnd("mp2", "mr2", nil),
							model.NewResourceEventEnd("mp2", "mr3", nil),
						},
					},
				},
			},
		},
	}
	ctx := context.Background()
	datastores, _ := newDatastores(t, ctx)
	for _, ds := range datastores {
		for _, test := range tests {
			name := fmt.Sprintf("%T-%s", ds, test.name)
			t.Run(name, func(t *testing.T) {
				for _, event := range test.args.events {
					require.NoError(t, ds.WriteEvent(ctx, event))
				}
				es, err := ds.EngineStatus(ctx)
				require.NoError(t, err)
				testingutil.AssertEvent(t, es, test.expected)
			})
		}
	}
}

func testQuery(t *testing.T, ctx context.Context, ds Datastore, fieldName string, fieldValue string, expected ...*model.Resource) {
	_testQuery(t, ctx, ds, fieldName, fieldValue, false, expected...)
}

func testQueryNoResult(t *testing.T, ctx context.Context, ds Datastore, fieldName string, fieldValue string) {
	_testQuery(t, ctx, ds, fieldName, fieldValue, false)
}

func testQueryUnrecognizedKey(t *testing.T, ctx context.Context, ds Datastore, fieldName string, fieldValue string) {
	_testQuery(t, ctx, ds, fieldName, fieldValue, true)
}

func _testQuery(t *testing.T, ctx context.Context, ds Datastore, fieldName string, fieldValue string, unrecognizedKey bool, expected ...*model.Resource) {
	query := fmt.Sprintf(`{
  "filter":{
    "%v": "%v"
  }
}`, fieldName, fieldValue)
	resourcesRead, err := ds.GetResources(ctx, []byte(query))
	if unrecognizedKey {
		require.ErrorContains(t, err, "unrecognized key")
	} else {
		require.NoError(t, err)
		testingutil.AssertEqualsResources(t, model.Resources(expected), resourcesRead.Resources)
	}
}

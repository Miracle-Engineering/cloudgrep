package datastore

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"testing"

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
			//check number of groups
			assert.Equal(t, 2, len(fields))
			//check fields by group
			assert.Equal(t, 2, len(fields.FindGroup("core").Fields))
			assert.Equal(t, 9, len(fields.FindGroup("tags").Fields))

			//test a few fields
			model.AssertEqualsField(t, model.Field{
				Name:  "region",
				Count: 3,
				Values: model.FieldValues{
					model.FieldValue{Value: "us-east-1", Count: 3},
				}}, *fields.FindField("core", "region"))

			typeField := *fields.FindField("core", "type")
			model.AssertEqualsField(t, model.Field{
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
				Name:  "team",
				Count: 2,
				Values: model.FieldValues{
					model.FieldValue{Value: "infra", Count: 1},
					model.FieldValue{Value: "dev", Count: 1},
				}}, *fields.FindField("tags", "team"))

			//test long field
			model.AssertEqualsField(t, model.Field{
				Name:  tagMaxKey,
				Count: 1,
				Values: model.FieldValues{
					model.FieldValue{Value: tagMaxValue, Count: 1},
				}}, *fields.FindField("tags", tagMaxKey))

			//test the tag field called "region"
			model.AssertEqualsField(t, model.Field{
				Name:  "region",
				Count: 1,
				Values: model.FieldValues{
					model.FieldValue{Value: "us-west-2", Count: 1},
				}}, *fields.FindField("tags", "region"))

		})
	}
}

func TestEngineStatus(t *testing.T) {
	engineStatuses := testdata.GetEngineStatus(t)
	ctx := context.Background()
	for _, datastore := range newDatastores(t, ctx) {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {
			err := datastore.WriteEngineStatusStart(ctx, "mock")
			if err != nil && err.Error() == "not implemented" {
				return
			}

			status, err := datastore.GetEngineStatus(ctx)
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}
			//check stats
			assert.NoError(t, err)
			model.AssertEqualsEngineStatus(t, engineStatuses[0], status)

			err = datastore.WriteEngineStatusEnd(ctx, "mock", nil)
			if err != nil && err.Error() == "not implemented" {
				return
			}

			status, err = datastore.GetEngineStatus(ctx)
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}
			//check stats
			assert.NoError(t, err)
			model.AssertEqualsEngineStatus(t, engineStatuses[1], status)

			err = datastore.WriteEngineStatusStart(ctx, "mock")
			if err != nil && err.Error() == "not implemented" {
				return
			}

			status, err = datastore.GetEngineStatus(ctx)
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}
			//check stats
			assert.NoError(t, err)
			model.AssertEqualsEngineStatus(t, engineStatuses[0], status)

			err = datastore.WriteEngineStatusEnd(ctx, "mock", errors.New(engineStatuses[2].ErrorMessage))
			if err != nil && err.Error() == "not implemented" {
				return
			}

			status, err = datastore.GetEngineStatus(ctx)
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}
			//check stats
			assert.NoError(t, err)
			model.AssertEqualsEngineStatus(t, engineStatuses[2], status)

		})
	}
}

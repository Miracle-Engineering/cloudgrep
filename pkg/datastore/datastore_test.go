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
			assert.NotZero(t, len(resources))

			assert.NoError(t, datastore.WriteResources(ctx, resources))

			var resourcesRead []*model.Resource
			resourcesRead, err := datastore.GetResources(ctx, model.EmptyFilter())
			assert.NoError(t, err)
			assert.Equal(t, len(resources), len(resourcesRead))
			util.AssertEqualsResources(t, resources, resourcesRead)

			//test getting a specific resource
			for _, r := range resources {
				resource, err := datastore.GetResource(ctx, r.Id)
				assert.NoError(t, err)
				util.AssertEqualsResourcePter(t, r, resource)
			}

		})
	}
}
func TestResourcesFiltering(t *testing.T) {
	ctx := context.Background()
	for _, datastore := range newDatastores(t, ctx) {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			all_resources := testdata.GetResources(t)
			resourceInst1 := all_resources[0]  //team:infra, release tag
			resourceInst2 := all_resources[1]  //team:dev, no release tag
			resourceBucket := all_resources[2] //s3 bucket without tags
			resourceLB := all_resources[3]     //a load balancer

			assert.NoError(t, datastore.WriteResources(ctx, all_resources))

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
			util.AssertEqualsResourcePter(t, resourceInst1, resourcesRead[0])

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
			util.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead)

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
			util.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead)

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
			util.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, resourcesRead)

			//test exclude - returns the resources without the tag release
			filter = model.Filter{
				Tags: []model.Tag{
					{Key: "release", Value: "*", Exclude: true},
				},
			}
			resourcesRead, err = datastore.GetResources(ctx, filter)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(resourcesRead))
			util.AssertEqualsResources(t, model.Resources{resourceInst2, resourceBucket, resourceLB}, resourcesRead)

			//test 2 exclusions - each instance resource has 1 tag but not both, kept them
			// TODO we should probably revisit this behavior -> exclude tags should be an OR not AND
			filter = model.Filter{
				Tags: []model.Tag{
					{Key: "release", Value: "*", Exclude: true},
					{Key: "debug:info", Value: "*", Exclude: true},
				},
			}
			resourcesRead, err = datastore.GetResources(ctx, filter)
			assert.NoError(t, err)
			util.AssertEqualsResources(t, all_resources, resourcesRead)

			//mix include and exclude filters
			filter = model.Filter{
				Tags: []model.Tag{
					{Key: "release", Value: "*", Exclude: true},
					{Key: "vpc", Value: "vpc-123"},
				},
			}
			resourcesRead, err = datastore.GetResources(ctx, filter)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(resourcesRead))
			util.AssertEqualsResourcePter(t, resourceInst2, resourcesRead[0])

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
			assert.Equal(t, model.Stats{ResourcesCount: 4}, stats)

		})
	}
}

func TestTagsFiltering(t *testing.T) {
	ctx := context.Background()
	for _, datastore := range newDatastores(t, ctx) {
		name := fmt.Sprintf("%T", datastore)
		t.Run(name, func(t *testing.T) {

			all_resources := testdata.GetResources(t)
			assert.NoError(t, datastore.WriteResources(ctx, all_resources))

			var tags model.TagInfos

			tags, err := datastore.GetTags(ctx, model.EmptyFilter(), 10)
			//do not test result if not implemented
			if err != nil && err.Error() == "not implemented" {
				return
			}

			assert.NoError(t, err)
			//check all tags returned
			assert.Equal(t, 10, len(tags))
			//check a specific entry
			clusterTag := model.TagInfo{
				Key:         "cluster",
				Values:      []string{"prod-cluster", "dev-cluster"},
				ResourceIds: []model.ResourceId{"i-123", "i-124", "arn:aws:elasticloadbalancing:us-east-1:248233625043:loadbalancer/net/opta-staging-ingress/14522ba1bd959dd6"},
				Count:       3,
			}
			util.AssertEqualsTagInfo(t, clusterTag, *tags.Find("cluster"))

			//test limit - since it's sorted, the most frequent tag is always returned first
			tags, err = datastore.GetTags(ctx, model.EmptyFilter(), 1)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(tags))
			util.AssertEqualsTagInfo(t, clusterTag, *tags[0])

			//test filter is applied
			filter := model.Filter{
				Tags: []model.Tag{{Key: "service.k8s.aws/resource", Value: "LoadBalancer"}},
			}
			tags, err = datastore.GetTags(ctx, filter, 10)
			assert.NoError(t, err)
			//only the lb should be scoped
			//the lb has 3 tags
			assert.Equal(t, 3, len(tags))
			//only he values and resource related to the LB are present
			clusterTag = model.TagInfo{
				Key:         "cluster",
				Values:      []string{"prod-cluster"},
				ResourceIds: []model.ResourceId{"arn:aws:elasticloadbalancing:us-east-1:248233625043:loadbalancer/net/opta-staging-ingress/14522ba1bd959dd6"},
				Count:       1,
			}
			util.AssertEqualsTagInfo(t, clusterTag, *tags.Find("cluster"))

			//test exclude filter - all resource have this tag, no result expected
			filter = model.Filter{
				Tags: []model.Tag{{Key: "cluster", Value: "*", Exclude: true}},
			}
			tags, err = datastore.GetTags(ctx, filter, 10)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(tags))

		})
	}
}

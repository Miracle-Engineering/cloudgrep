package aws

import (
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	testprovider "github.com/juandiegopalomino/cloudgrep/pkg/testingutil/provider"
)

func TestFetchElasticacheClusters(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "elasticache.CacheCluster")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:            "elasticache.CacheCluster",
		Region:          defaultRegion,
		DisplayIdPrefix: "test-0-",
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "elasticache-cluster-0",
			},
		},
	})
}

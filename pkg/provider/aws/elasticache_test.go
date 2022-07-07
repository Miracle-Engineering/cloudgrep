package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchElasticacheClusters(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "elasticache.CacheCluster")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		AccountId: ctx.accountId,
		Type:      "elasticache.CacheCluster",
		Region:    defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "elasticache-cluster-0",
			},
		},
	})
}

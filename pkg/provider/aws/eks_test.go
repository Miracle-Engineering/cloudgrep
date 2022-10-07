package aws

import (
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	testprovider "github.com/juandiegopalomino/cloudgrep/pkg/testingutil/provider"
)

func TestFetchEKSClusters(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "eks.Cluster")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:            "eks.Cluster",
		Region:          defaultRegion,
		DisplayIdPrefix: "main",
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "eks-cluster-main",
			},
		},
	})
}

func TestFetchEKSNodegroup(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "eks.Nodegroup")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:            "eks.Nodegroup",
		Region:          defaultRegion,
		DisplayIdPrefix: "main-default-",
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "eks-cluster-main-default-node-group",
			},
		},
	})
}

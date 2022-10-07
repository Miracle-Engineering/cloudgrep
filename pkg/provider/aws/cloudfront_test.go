package aws

import (
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	testprovider "github.com/juandiegopalomino/cloudgrep/pkg/testingutil/provider"
)

func TestFetchCloudfrontDistributions(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "cloudfront.Distribution")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "cloudfront.Distribution",
		Region: globalRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "cloudfront-distribution-0",
			},
		},
	})
}

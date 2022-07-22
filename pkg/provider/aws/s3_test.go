package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchS3Buckets(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "s3.Bucket")

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:            "s3.Bucket",
		Region:          "global",
		DisplayIdPrefix: "cloudgrep-testing-0-",
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "s3-bucket-0",
			},
		},
	})
}

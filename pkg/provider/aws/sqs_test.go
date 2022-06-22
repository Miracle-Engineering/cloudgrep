package aws

import (
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
	"testing"
)

func TestFetchSqsQueue(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "sqs.SQS")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "sqs.SQS",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "sqs-queue-0",
			},
		},
	})
}

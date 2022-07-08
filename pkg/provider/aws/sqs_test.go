package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchSqsQueue(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "sqs.Queue")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:            "sqs.Queue",
		DisplayIdPrefix: "testing-",
		Region:          defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "sqs-queue-0",
			},
		},
	})
}

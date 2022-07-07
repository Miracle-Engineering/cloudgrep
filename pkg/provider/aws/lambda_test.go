package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchFunctions(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "lambda.Function")

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		AccountId: ctx.accountId,
		Type:      "lambda.Function",
		Region:    defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "lambda-function-0",
			},
		},
		RawData: map[string]any{
			"Runtime": "nodejs16.x",
		},
	})
}

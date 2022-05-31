package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFetchFunctions(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	raw := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchFunctions)
	resources := testingutil.ConvertToResources(t, ctx.ctx, ctx.p.mapper, raw)

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceCount(t, resources, "lambda-function-0", 1)
}

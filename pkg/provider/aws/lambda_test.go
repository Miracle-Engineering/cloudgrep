package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFetchFunctions(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchLambdaFunctions)

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceCount(t, resources, "lambda-function-0", 1)
}

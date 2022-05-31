package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFetchRDSInstances(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	instances := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchRDSInstances)
	resources := testingutil.ConvertToResources(t, ctx.ctx, ctx.p.mapper, instances)

	testingutil.AssertResourceCount(t, resources, "rds-instance-0", 1)
}

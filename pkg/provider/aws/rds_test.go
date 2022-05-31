package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFetchRDSInstances(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	raw := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchRDSInstances)
	resources := testingutil.ConvertToResources(t, ctx.ctx, ctx.p.mapper, raw)

	testingutil.AssertResourceCount(t, resources, "rds-instance-0", 1)
}

func TestFetchRDSClusters(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	raw := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchRDSClusters)
	resources := testingutil.ConvertToResources(t, ctx.ctx, ctx.p.mapper, raw)

	testingutil.AssertResourceCount(t, resources, "rds-cluster-0", 1)
}

package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchRDSInstances(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "rds.DBInstance")

	testingutil.AssertResourceCount(t, resources, "rds-instance-0", 1)
}

func TestFetchRDSClusters(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "rds.DBCluster")

	testingutil.AssertResourceCount(t, resources, "rds-cluster-0", 1)
}

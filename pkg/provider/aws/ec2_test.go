package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestEc2FetchInstances(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	instances := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchEC2Instances)
	resources := testingutil.ConvertToResources(t, ctx.ctx, ctx.p.mapper, instances)

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceCount(t, resources, "ec2-instance-0", 1)
}

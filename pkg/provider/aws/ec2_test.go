package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFetchEC2Instances(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchEC2Instances)

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceCount(t, resources, "ec2-instance-0", 1)
}

func TestFetchEBSVolumes(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchEBSVolumes)

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceCount(t, resources, "ec2-volume-0", 1)
}

package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFetchLoadBalancers(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	raw := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchLoadBalancers)
	resources := testingutil.ConvertToResources(t, ctx.ctx, ctx.p.mapper, raw)

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceCount(t, resources, "elb-alb-0", 1)
}

package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFetchS3Buckets(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	raw := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchS3Buckets)
	resources := testingutil.ConvertToResources(t, ctx.ctx, ctx.p.mapper, raw)

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceCount(t, resources, "s3-bucket-0", 1)
}

package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/testingutil"
)

func TestFetchS3Buckets(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testingutil.MustFetchAll(ctx.ctx, t, ctx.p.FetchS3Buckets)

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceCount(t, resources, "s3-bucket-0", 1)
}

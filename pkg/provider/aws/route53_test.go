package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchHealthChecks(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "route53.HealthCheck")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "route53.HealthCheck",
		Region: globalRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "route53-health-check-0",
			},
		},
	})
}

func TestFetchHostedZones(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "route53.HostedZone")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "route53.HostedZone",
		Region: globalRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "route53-hosted-zone-0",
			},
		},
		RawData: map[string]any{
			"Name": "0.example.com",
		},
	})
}

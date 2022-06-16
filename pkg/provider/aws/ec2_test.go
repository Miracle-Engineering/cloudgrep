package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchEC2Instances(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.Instance")

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.Instance",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-instance-0",
			},
		},
		RawData: map[string]any{
			"ImageId": "ami-0e449176cecc3e577",
		},
	})
}

func TestFetchEBSVolumes(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.Volume")

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.Volume",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-volume-0",
			},
		},
		RawData: map[string]any{
			"VolumeType": "gp2",
		},
	})
}

func TestFetchAddresses(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.Address")
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.Address",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-ip-0",
			},
		},
	})
}

func TestFetchImages(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.Image")
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.Image",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-ami-0",
			},
		},
		RawData: map[string]any{
			"Architecture": "x86_84",
		},
	})
}

func TestFetchKeyPairs(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.KeyPair")
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.KeyPair",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-keypair-0",
			},
		},
	})
}

func TestFetchLaunchTemplates(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.LaunchTemplate")

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.LaunchTemplate",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-launch-template-0",
			},
		},
	})
}

func TestFetchNatGateways(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.NatGateway")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.NatGateway",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "vpc-main-nat-main",
			},
		},
	})
}

func TestFetchNetworkAcl(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.NetworkAcl")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.NetworkAcl",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "vpc-main-acl-0",
			},
		},
	})
}

func TestFetchNetworkInterface(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.NetworkInterface")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.NetworkInterface",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-eni-0",
			},
		},
	})
}

func TestFetchRouteTable(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.RouteTable")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.RouteTable",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "vpc-main-route-table-private",
			},
		},
	})
}

func TestFetchSecurityGroups(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.SecurityGroup")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.SecurityGroup",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-sg-0",
			},
		},
	})
}

func TestFetchSnapshots(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.Snapshot")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.Snapshot",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "ec2-ebs-snapshot-0",
			},
		},
	})
}

func TestFetchSubnets(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.Subnet")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.Subnet",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "vpc-main-subnet-a",
			},
		},
	})
}

func TestFetchVpcs(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)
	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "ec2.Vpc")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "ec2.Vpc",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "vpc-main",
			},
		},
	})
}

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) registerEc2(mapping map[string]mapper) {
	mapping["ec2.Address"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2Address,
		IdField:           "AllocationId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.CapacityReservation"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2CapacityReservation,
		IdField:           "CapacityReservationId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.ClientVpnEndpoint"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2ClientVpnEndpoint,
		IdField:           "ClientVpnEndpointId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Fleet"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2Fleet,
		IdField:           "FleetId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.FlowLogs"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2FlowLogs,
		IdField:           "FlowLogId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Image"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2Image,
		IdField:           "ImageId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Instance"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2Instance,
		IdField:           "InstanceId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.KeyPair"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2KeyPair,
		IdField:           "KeyPairId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.LaunchTemplate"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2LaunchTemplate,
		IdField:           "LaunchTemplateId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NatGateway"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2NatGateway,
		IdField:           "NatGatewayId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NetworkAcl"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2NetworkAcl,
		IdField:           "NetworkAclId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NetworkInterface"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2NetworkInterface,
		IdField:           "NetworkInterfaceId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagSet",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.ReservedInstance"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2ReservedInstance,
		IdField:           "ReservedInstancesId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.RouteTable"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2RouteTable,
		IdField:           "RouteTableId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.SecurityGroup"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2SecurityGroup,
		IdField:           "GroupId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Snapshot"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2Snapshot,
		IdField:           "SnapshotId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.SpotInstanceRequest"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2SpotInstanceRequest,
		IdField:           "SpotInstanceRequestId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Subnet"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2Subnet,
		IdField:           "SubnetId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Volume"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2Volume,
		IdField:           "VolumeId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Vpc"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetchEc2Vpc,
		IdField:           "VpcId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetchEc2Address(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeAddressesInput{}

	resourceConverter := p.converterFor("ec2.Address")
	results, err := client.DescribeAddresses(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.Address", err)
	}
	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, results.Addresses); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetchEc2CapacityReservation(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeCapacityReservationsInput{}
	input.Filters = describeCapacityReservationsFilters()

	resourceConverter := p.converterFor("ec2.CapacityReservation")
	paginator := ec2.NewDescribeCapacityReservationsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.CapacityReservation", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.CapacityReservations); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2ClientVpnEndpoint(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeClientVpnEndpointsInput{}

	resourceConverter := p.converterFor("ec2.ClientVpnEndpoint")
	paginator := ec2.NewDescribeClientVpnEndpointsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.ClientVpnEndpoint", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.ClientVpnEndpoints); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2Fleet(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeFleetsInput{}
	input.Filters = describeFleetsFilters()

	resourceConverter := p.converterFor("ec2.Fleet")
	paginator := ec2.NewDescribeFleetsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Fleet", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Fleets); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2FlowLogs(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeFlowLogsInput{}

	resourceConverter := p.converterFor("ec2.FlowLogs")
	paginator := ec2.NewDescribeFlowLogsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.FlowLogs", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.FlowLogs); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2Image(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeImagesInput{}
	input.Owners = describeImagesOwners()

	resourceConverter := p.converterFor("ec2.Image")
	results, err := client.DescribeImages(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.Image", err)
	}
	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, results.Images); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetchEc2Instance(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeInstancesInput{}
	input.Filters = describeInstancesFilters()

	resourceConverter := p.converterFor("ec2.Instance")
	paginator := ec2.NewDescribeInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Instance", err)
		}

		for _, item_0 := range page.Reservations {
			if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, item_0.Instances); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Provider) fetchEc2KeyPair(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeKeyPairsInput{}

	resourceConverter := p.converterFor("ec2.KeyPair")
	results, err := client.DescribeKeyPairs(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.KeyPair", err)
	}
	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, results.KeyPairs); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetchEc2LaunchTemplate(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeLaunchTemplatesInput{}

	resourceConverter := p.converterFor("ec2.LaunchTemplate")
	paginator := ec2.NewDescribeLaunchTemplatesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.LaunchTemplate", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.LaunchTemplates); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2NatGateway(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeNatGatewaysInput{}

	resourceConverter := p.converterFor("ec2.NatGateway")
	paginator := ec2.NewDescribeNatGatewaysPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.NatGateway", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.NatGateways); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2NetworkAcl(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeNetworkAclsInput{}

	resourceConverter := p.converterFor("ec2.NetworkAcl")
	paginator := ec2.NewDescribeNetworkAclsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.NetworkAcl", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.NetworkAcls); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2NetworkInterface(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeNetworkInterfacesInput{}

	resourceConverter := p.converterFor("ec2.NetworkInterface")
	paginator := ec2.NewDescribeNetworkInterfacesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.NetworkInterface", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.NetworkInterfaces); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2ReservedInstance(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeReservedInstancesInput{}
	input.Filters = describeReservedInstancesFilters()

	resourceConverter := p.converterFor("ec2.ReservedInstance")
	results, err := client.DescribeReservedInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.ReservedInstance", err)
	}
	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, results.ReservedInstances); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetchEc2RouteTable(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeRouteTablesInput{}

	resourceConverter := p.converterFor("ec2.RouteTable")
	paginator := ec2.NewDescribeRouteTablesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.RouteTable", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.RouteTables); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2SecurityGroup(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeSecurityGroupsInput{}

	resourceConverter := p.converterFor("ec2.SecurityGroup")
	paginator := ec2.NewDescribeSecurityGroupsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.SecurityGroup", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.SecurityGroups); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2Snapshot(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeSnapshotsInput{}
	input.OwnerIds = describeSnapshotsOwners()

	resourceConverter := p.converterFor("ec2.Snapshot")
	paginator := ec2.NewDescribeSnapshotsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Snapshot", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Snapshots); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2SpotInstanceRequest(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeSpotInstanceRequestsInput{}
	input.Filters = describeSpotInstanceRequestsFilters()

	resourceConverter := p.converterFor("ec2.SpotInstanceRequest")
	paginator := ec2.NewDescribeSpotInstanceRequestsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.SpotInstanceRequest", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.SpotInstanceRequests); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2Subnet(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeSubnetsInput{}

	resourceConverter := p.converterFor("ec2.Subnet")
	paginator := ec2.NewDescribeSubnetsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Subnet", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Subnets); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2Volume(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeVolumesInput{}

	resourceConverter := p.converterFor("ec2.Volume")
	paginator := ec2.NewDescribeVolumesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Volume", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Volumes); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchEc2Vpc(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeVpcsInput{}

	resourceConverter := p.converterFor("ec2.Vpc")
	paginator := ec2.NewDescribeVpcsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Vpc", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Vpcs); err != nil {
			return err
		}
	}

	return nil
}

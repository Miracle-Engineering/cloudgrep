package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_ec2(mapping map[string]mapper) {
	mapping["ec2.Address"] = mapper{
		FetchFunc: p.fetch_ec2_Address,
		IdField:   "AllocationId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.CapacityReservation"] = mapper{
		FetchFunc: p.fetch_ec2_CapacityReservation,
		IdField:   "CapacityReservationId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.ClientVpnEndpoint"] = mapper{
		FetchFunc: p.fetch_ec2_ClientVpnEndpoint,
		IdField:   "ClientVpnEndpointId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.ElasticGpu"] = mapper{
		FetchFunc: p.fetch_ec2_ElasticGpu,
		IdField:   "ElasticGpuId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Fleet"] = mapper{
		FetchFunc: p.fetch_ec2_Fleet,
		IdField:   "FleetId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.FlowLogs"] = mapper{
		FetchFunc: p.fetch_ec2_FlowLogs,
		IdField:   "FlowLogId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Image"] = mapper{
		FetchFunc: p.fetch_ec2_Image,
		IdField:   "ImageId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Instance"] = mapper{
		FetchFunc: p.fetch_ec2_Instance,
		IdField:   "InstanceId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.KeyPair"] = mapper{
		FetchFunc: p.fetch_ec2_KeyPair,
		IdField:   "KeyPairId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.LaunchTemplate"] = mapper{
		FetchFunc: p.fetch_ec2_LaunchTemplate,
		IdField:   "LaunchTemplateId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NatGateway"] = mapper{
		FetchFunc: p.fetch_ec2_NatGateway,
		IdField:   "NatGatewayId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NetworkAcl"] = mapper{
		FetchFunc: p.fetch_ec2_NetworkAcl,
		IdField:   "NetworkAclId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NetworkInterface"] = mapper{
		FetchFunc: p.fetch_ec2_NetworkInterface,
		IdField:   "NetworkInterfaceId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "TagSet",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.ReservedInstance"] = mapper{
		FetchFunc: p.fetch_ec2_ReservedInstance,
		IdField:   "ReservedInstancesId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.RouteTable"] = mapper{
		FetchFunc: p.fetch_ec2_RouteTable,
		IdField:   "RouteTableId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.SecurityGroup"] = mapper{
		FetchFunc: p.fetch_ec2_SecurityGroup,
		IdField:   "GroupId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Snapshot"] = mapper{
		FetchFunc: p.fetch_ec2_Snapshot,
		IdField:   "SnapshotId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.SpotInstanceRequest"] = mapper{
		FetchFunc: p.fetch_ec2_SpotInstanceRequest,
		IdField:   "SpotInstanceRequestId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Subnet"] = mapper{
		FetchFunc: p.fetch_ec2_Subnet,
		IdField:   "SubnetId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Volume"] = mapper{
		FetchFunc: p.fetch_ec2_Volume,
		IdField:   "VolumeId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Vpc"] = mapper{
		FetchFunc: p.fetch_ec2_Vpc,
		IdField:   "VpcId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetch_ec2_Address(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_CapacityReservation(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_ClientVpnEndpoint(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_ElasticGpu(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeElasticGpusInput{}

	resourceConverter := p.converterFor("ec2.ElasticGpu")
	results, err := client.DescribeElasticGpus(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.ElasticGpu", err)
	}
	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, results.ElasticGpuSet); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetch_ec2_Fleet(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_FlowLogs(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_Image(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_Instance(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_KeyPair(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_LaunchTemplate(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_NatGateway(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_NetworkAcl(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_NetworkInterface(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_ReservedInstance(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_RouteTable(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_SecurityGroup(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_Snapshot(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_SpotInstanceRequest(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_Subnet(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_Volume(ctx context.Context, output chan<- model.Resource) error {
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

func (p *Provider) fetch_ec2_Vpc(ctx context.Context, output chan<- model.Resource) error {
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

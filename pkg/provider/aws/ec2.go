package aws

import (
	"context"
	"fmt"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/run-x/cloudgrep/pkg/model"
)

func (p *Provider) FetchEC2Instances(ctx context.Context, output chan<- model.Resource) error {
	resourceType := "ec2.Instance"
	ec2Client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeInstancesInput{}
	paginator := ec2.NewDescribeInstancesPaginator(ec2Client, input)

	resourceConverter := p.converterFor(resourceType)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 Instances: %w", err)
		}

		for _, r := range page.Reservations {
			if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, r.Instances); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Provider) FetchEBSVolumes(ctx context.Context, output chan<- model.Resource) error {
	resourceType := "ec2.Volume"
	ec2Client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeVolumesInput{}
	paginator := ec2.NewDescribeVolumesPaginator(ec2Client, input)

	resourceConverter := p.converterFor(resourceType)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 EBS Volumes: %w", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Volumes); err != nil {
			return err
		}
	}

	return nil
}

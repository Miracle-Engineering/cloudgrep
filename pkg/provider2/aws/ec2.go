package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
	"github.com/run-x/cloudgrep/pkg/util"
)

func (p *Provider) FetchEC2Instances(ctx context.Context, output chan<- model.Resource) error {
	resourceType := "ec2.Instance"
	ec2Client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeInstancesInput{}
	paginator := ec2.NewDescribeInstancesPaginator(ec2Client, input)
	mapping := p.getTypeMapping()[resourceType]

	resourceConverter := &resourceconverter.ReflectionConverter{
		Region:       p.config.Region,
		ResourceType: resourceType,
		TagField:     mapping.TagField,
		IdField:      mapping.IdField,
	}
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 Instances: %w", err)
		}

		for _, r := range page.Reservations {
			if err := SendAllConverted(ctx, output, resourceConverter, r.Instances); err != nil {
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
	mapping := p.getTypeMapping()[resourceType]

	resourceConverter := resourceconverter.ReflectionConverter{
		Region:       p.config.Region,
		ResourceType: resourceType,
		TagField:     mapping.TagField,
		IdField:      mapping.IdField,
	}
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 EBS Volumes: %w", err)
		}

		var resources []model.Resource
		for _, i := range page.Volumes {
			newResource, err := resourceConverter.ToResource(ctx, i, nil)
			if err != nil {
				return err
			}
			resources = append(resources, newResource)
		}

		if err := util.SendAllFromSlice(ctx, output, resources); err != nil {
			return err
		}
	}

	return nil
}

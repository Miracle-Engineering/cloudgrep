package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (awsPrv *AWSProvider) FetchEC2Instances(ctx context.Context, output chan<- types.Instance) error {
	input := &ec2.DescribeInstancesInput{}
	p := ec2.NewDescribeInstancesPaginator(awsPrv.ec2Client, input)
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 Instances: %w", err)
		}

		for _, r := range page.Reservations {
			for _, i := range r.Instances {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case output <- i:
				}
			}
		}
	}

	return nil
}

func (p *AWSProvider) FetchEBSVolumes(ctx context.Context) ([]types.Volume, error) {
	input := &ec2.DescribeVolumesInput{}

	result, err := p.ec2Client.DescribeVolumes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch EC2 EBS Volumes: %w", err)
	}
	return result.Volumes, nil
}

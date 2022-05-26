package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/run-x/cloudgrep/pkg/util"
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
			if err := util.SendAllFromSlice(ctx, output, r.Instances); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *AWSProvider) FetchEBSVolumes(ctx context.Context, output chan<- types.Volume) error {
	input := &ec2.DescribeVolumesInput{}
	paginator := ec2.NewDescribeVolumesPaginator(p.ec2Client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 EBS Volumes: %w", err)
		}

		if err := util.SendAllFromSlice(ctx, output, page.Volumes); err != nil {
			return err
		}
	}

	return nil
}

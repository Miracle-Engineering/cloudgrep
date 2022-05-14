package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (awsPrv *AWSProvider) FetchEC2Instances(ctx context.Context) ([]types.Instance, error) {
	input := &ec2.DescribeInstancesInput{}
	var instances []types.Instance
	//TODO use pagination (consider returning a channel?)
	result, err := awsPrv.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch EC2 Instances: %w", err)
	}

	for _, r := range result.Reservations {
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}

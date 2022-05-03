package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (awsPrv *AWSProvider) FetchEC2Instances(ctx context.Context) ([]types.Instance, error) {
	input := &ec2.DescribeInstancesInput{}
	var instances []types.Instance
	//TODO use pagination (consider returning a channel?)
	result, err := awsPrv.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return []types.Instance{}, err
	}

	for _, r := range result.Reservations {
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}

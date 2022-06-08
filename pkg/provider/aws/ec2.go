package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func describeInstancesFilters() []types.Filter {
	return []types.Filter{
		{
			Name: aws.String("instance-state-name"),
			Values: []string{
				"pending",
				"running",
				"shutting-down",
				"stopped",
				"stopping",
			},
		},
	}
}

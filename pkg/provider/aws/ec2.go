package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func describeCapacityReservationsFilters() []types.Filter {
	return []types.Filter{
		{
			Name: aws.String("state"),
			Values: []string{
				"active",
				"pending",
			},
		},
	}
}

func describeFleetsFilters() []types.Filter {
	return []types.Filter{
		{
			Name: aws.String("fleet-state"),
			Values: []string{
				"active",
				"deleted-running",
				"deleted-terminating",
				"modifying",
				"submitted",
			},
		},
	}
}

func describeImagesOwners() []string {
	return []string{"self"}
}

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

func describeReservedInstancesFilters() []types.Filter {
	return []types.Filter{
		{
			Name: aws.String("state"),
			Values: []string{
				"payment-pending",
				"active",
				"payment-failed",
			},
		},
	}
}

func describeSnapshotsOwners() []string {
	return []string{"self"}
}

func describeSpotInstanceRequestsFilters() []types.Filter {
	return []types.Filter{
		{
			Name: aws.String("state"),
			Values: []string{
				"open",
				"active",
			},
		},
	}
}

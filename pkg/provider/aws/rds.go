package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func (p *AWSProvider) FetchRDSInstances(ctx context.Context) ([]types.DBInstance, error) {
	input := &rds.DescribeDBInstancesInput{}
	paginator := rds.NewDescribeDBInstancesPaginator(p.rdsClient, input)

	var resources []types.DBInstance

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch RDS DB Instances: %w", err)
		}

		resources = append(resources, output.DBInstances...)
	}

	return resources, nil
}

func (p *AWSProvider) FetchRDSClusters(ctx context.Context) ([]types.DBCluster, error) {
	input := &rds.DescribeDBClustersInput{}
	paginator := rds.NewDescribeDBClustersPaginator(p.rdsClient, input)

	var resources []types.DBCluster

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch RDS DB Clusters: %w", err)
		}

		resources = append(resources, output.DBClusters...)
	}

	return resources, nil
}

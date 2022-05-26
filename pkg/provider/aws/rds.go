package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/run-x/cloudgrep/pkg/util"
)

func (p *AWSProvider) FetchRDSInstances(ctx context.Context, output chan<- types.DBInstance) error {
	input := &rds.DescribeDBInstancesInput{}
	paginator := rds.NewDescribeDBInstancesPaginator(p.rdsClient, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch RDS DB Instances: %w", err)
		}

		if err := util.SendAllFromSlice(ctx, output, page.DBInstances); err != nil {
			return err
		}
	}

	return nil
}

func (p *AWSProvider) FetchRDSClusters(ctx context.Context, output chan<- types.DBCluster) error {
	input := &rds.DescribeDBClustersInput{}
	paginator := rds.NewDescribeDBClustersPaginator(p.rdsClient, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch RDS DB Clusters: %w", err)
		}

		if err := util.SendAllFromSlice(ctx, output, page.DBClusters); err != nil {
			return err
		}
	}

	return nil
}

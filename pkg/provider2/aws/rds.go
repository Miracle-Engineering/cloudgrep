package aws

import (
	"context"
	"fmt"
	"github.com/run-x/cloudgrep/pkg/model"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func (p *Provider) FetchRDSInstances(ctx context.Context, output chan<- model.Resource) error {
	resourceType := "rds.DBInstance"
	rdsClient := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBInstancesInput{}
	paginator := rds.NewDescribeDBInstancesPaginator(rdsClient, input)

	resourceConverter := p.converterFor(resourceType)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch RDS Instances: %w", err)
		}

		if err := SendAllConverted(ctx, output, resourceConverter, page.DBInstances); err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) FetchRDSClusters(ctx context.Context, output chan<- model.Resource) error {
	resourceType := "rds.DBCluster"
	rdsClient := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBClustersInput{}
	paginator := rds.NewDescribeDBClustersPaginator(rdsClient, input)

	resourceConverter := p.converterFor(resourceType)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch RDS Clusters: %w", err)
		}

		if err := SendAllConverted(ctx, output, resourceConverter, page.DBClusters); err != nil {
			return err
		}
	}
	return nil
}

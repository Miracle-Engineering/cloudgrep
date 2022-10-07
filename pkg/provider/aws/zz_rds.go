package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) registerRds(mapping map[string]mapper) {
	mapping["rds.DBCluster"] = mapper{
		ServiceEndpointID: "rds",
		FetchFunc:         p.fetchRdsDBCluster,
		IdField:           "DBClusterIdentifier",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagList",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["rds.DBClusterSnapshot"] = mapper{
		ServiceEndpointID: "rds",
		FetchFunc:         p.fetchRdsDBClusterSnapshot,
		IdField:           "DBClusterSnapshotIdentifier",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagList",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["rds.DBInstance"] = mapper{
		ServiceEndpointID: "rds",
		FetchFunc:         p.fetchRdsDBInstance,
		IdField:           "DBInstanceIdentifier",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagList",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["rds.DBSnapshot"] = mapper{
		ServiceEndpointID: "rds",
		FetchFunc:         p.fetchRdsDBSnapshot,
		IdField:           "DBSnapshotIdentifier",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagList",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetchRdsDBCluster(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBClustersInput{}

	resourceConverter := p.converterFor("rds.DBCluster")
	paginator := rds.NewDescribeDBClustersPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBCluster", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.DBClusters); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchRdsDBClusterSnapshot(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBClusterSnapshotsInput{}

	resourceConverter := p.converterFor("rds.DBClusterSnapshot")
	paginator := rds.NewDescribeDBClusterSnapshotsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBClusterSnapshot", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.DBClusterSnapshots); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchRdsDBInstance(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBInstancesInput{}

	resourceConverter := p.converterFor("rds.DBInstance")
	paginator := rds.NewDescribeDBInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBInstance", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.DBInstances); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetchRdsDBSnapshot(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBSnapshotsInput{}

	resourceConverter := p.converterFor("rds.DBSnapshot")
	paginator := rds.NewDescribeDBSnapshotsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBSnapshot", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.DBSnapshots); err != nil {
			return err
		}
	}

	return nil
}

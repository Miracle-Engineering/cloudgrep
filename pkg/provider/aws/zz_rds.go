package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_rds(mapping map[string]mapper) {
	mapping["rds.DBCluster"] = mapper{
		ServiceEndpointID: "rds",
		FetchFunc:         p.fetch_rds_DBCluster,
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
		FetchFunc:         p.fetch_rds_DBClusterSnapshot,
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
		FetchFunc:         p.fetch_rds_DBInstance,
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
		FetchFunc:         p.fetch_rds_DBSnapshot,
		IdField:           "DBSnapshotIdentifier",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagList",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetch_rds_DBCluster(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBClustersInput{}

	commonTransformers := p.baseTransformers("rds.DBCluster")
	converter := p.converterFor("rds.DBCluster")
	transformers := append(
		resourceconverter.AllToGeneric[types.DBCluster](commonTransformers...),
		resourceconverter.WithConverter[types.DBCluster](converter),
	)
	paginator := rds.NewDescribeDBClustersPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBCluster", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.DBClusters, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_rds_DBClusterSnapshot(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBClusterSnapshotsInput{}

	commonTransformers := p.baseTransformers("rds.DBClusterSnapshot")
	converter := p.converterFor("rds.DBClusterSnapshot")
	transformers := append(
		resourceconverter.AllToGeneric[types.DBClusterSnapshot](commonTransformers...),
		resourceconverter.WithConverter[types.DBClusterSnapshot](converter),
	)
	paginator := rds.NewDescribeDBClusterSnapshotsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBClusterSnapshot", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.DBClusterSnapshots, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_rds_DBInstance(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBInstancesInput{}

	commonTransformers := p.baseTransformers("rds.DBInstance")
	converter := p.converterFor("rds.DBInstance")
	transformers := append(
		resourceconverter.AllToGeneric[types.DBInstance](commonTransformers...),
		resourceconverter.WithConverter[types.DBInstance](converter),
	)
	paginator := rds.NewDescribeDBInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBInstance", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.DBInstances, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_rds_DBSnapshot(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBSnapshotsInput{}

	commonTransformers := p.baseTransformers("rds.DBSnapshot")
	converter := p.converterFor("rds.DBSnapshot")
	transformers := append(
		resourceconverter.AllToGeneric[types.DBSnapshot](commonTransformers...),
		resourceconverter.WithConverter[types.DBSnapshot](converter),
	)
	paginator := rds.NewDescribeDBSnapshotsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBSnapshot", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.DBSnapshots, transformers...); err != nil {
			return err
		}
	}

	return nil
}

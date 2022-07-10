package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/hashicorp/go-multierror"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_eks(mapping map[string]mapper) {
	mapping["eks.Cluster"] = mapper{
		FetchFunc: p.fetch_eks_Cluster,
		IdField:   "Arn",
		IsGlobal:  false,
	}
	mapping["eks.Nodegroup"] = mapper{
		FetchFunc: p.fetch_eks_Nodegroup,
		IdField:   "NodegroupArn",
		IsGlobal:  false,
	}
}

func (p *Provider) get_cluster_names(ctx context.Context) ([]string, error) {
	client := eks.NewFromConfig(p.config)
	input := &eks.ListClustersInput{}

	paginator := eks.NewListClustersPaginator(client, input)
	var clusterNames []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s: %w", "eks.Cluster", err)
		}
		clusterNames = append(clusterNames, page.Clusters...)
	}
	return clusterNames, nil
}

func (p *Provider) fetch_eks_Cluster(ctx context.Context, output chan<- model.Resource) error {
	var err error
	client := eks.NewFromConfig(p.config)
	resourceConverter := p.converterFor("eks.Cluster")
	clusterNames, err := p.get_cluster_names(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "eks.Cluster", err)
	}
	var errors *multierror.Error
	var clusters []types.Cluster
	for _, clusterName := range clusterNames {
		describeClusterInput := &eks.DescribeClusterInput{
			Name: &clusterName,
		}
		describeClusterResults, err := client.DescribeCluster(ctx, describeClusterInput)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
		clusters = append(clusters, *describeClusterResults.Cluster)
	}

	var transformers resourceconverter.Transformers[types.Cluster]
	transformers.AddTags(p.getTags_eks_Cluster)

	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, clusters, transformers); err != nil {
		return err
	}

	return nil
}

func (p *Provider) getTags_eks_Cluster(ctx context.Context, resource types.Cluster) (model.Tags, error) {
	var tags model.Tags

	for key, value := range resource.Tags {
		tags = append(tags, model.Tag{
			Key:   key,
			Value: value,
		})
	}

	return tags, nil
}

func (p *Provider) fetch_eks_Nodegroup(ctx context.Context, output chan<- model.Resource) error {
	var err error
	var multiErrors *multierror.Error
	client := eks.NewFromConfig(p.config)
	clusterNames, err := p.get_cluster_names(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "eks.Nodegroup", err)
	}

	var transformers resourceconverter.Transformers[types.Nodegroup]
	transformers.AddTags(p.getTags_eks_Nodegroup)

	for _, clusterName := range clusterNames {
		input := &eks.ListNodegroupsInput{ClusterName: &clusterName}

		resourceConverter := p.converterFor("eks.Nodegroup")
		paginator := eks.NewListNodegroupsPaginator(client, input)
		var nodeGroupNames []string
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)

			if err != nil {
				return fmt.Errorf("failed to fetch %s: %w", "eks.Nodegroup", err)
			}
			nodeGroupNames = append(nodeGroupNames, page.Nodegroups...)
		}
		var nodegroups []types.Nodegroup
		for _, nodegroupName := range nodeGroupNames {
			describeNodegroupInput := &eks.DescribeNodegroupInput{
				ClusterName:   &clusterName,
				NodegroupName: &nodegroupName,
			}
			describeNodegroupResults, err := client.DescribeNodegroup(ctx, describeNodegroupInput)
			if err != nil {
				fmt.Println(err)
				multiErrors = multierror.Append(multiErrors, err)
			}
			nodegroups = append(nodegroups, *describeNodegroupResults.Nodegroup)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, nodegroups, transformers); err != nil {
			return err
		}
	}

	return multiErrors.ErrorOrNil()
}

func (p *Provider) getTags_eks_Nodegroup(ctx context.Context, resource types.Nodegroup) (model.Tags, error) {
	var tags model.Tags

	for key, value := range resource.Tags {
		tags = append(tags, model.Tag{
			Key:   key,
			Value: value,
		})
	}

	return tags, nil
}

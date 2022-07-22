package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) registerElb(mapping map[string]mapper) {
	mapping["elb.LoadBalancer"] = mapper{
		ServiceEndpointID: "elasticloadbalancing",
		FetchFunc:         p.fetchElbLoadBalancer,
		IdField:           "LoadBalancerArn",
		IsGlobal:          false,
	}
}

func (p *Provider) fetchElbLoadBalancer(ctx context.Context, output chan<- model.Resource) error {
	client := elasticloadbalancingv2.NewFromConfig(p.config)
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}

	resourceConverter := p.converterFor("elb.LoadBalancer")
	paginator := elasticloadbalancingv2.NewDescribeLoadBalancersPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "elb.LoadBalancer", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.LoadBalancers, p.getTagsElbLoadBalancer); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTagsElbLoadBalancer(ctx context.Context, resource types.LoadBalancer) (model.Tags, error) {
	client := elasticloadbalancingv2.NewFromConfig(p.config)
	input := &elasticloadbalancingv2.DescribeTagsInput{}

	input.ResourceArns = []string{*resource.LoadBalancerArn}

	output, err := client.DescribeTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "elb.LoadBalancer", err)
	}
	tagField_0 := output.TagDescriptions
	var tagField_1 []types.Tag
	for _, field := range tagField_0 {
		tagField_1 = append(tagField_1, field.Tags...)
	}

	var tags model.Tags

	for _, field := range tagField_1 {
		tags = append(tags, model.Tag{
			Key:   *field.Key,
			Value: *field.Value,
		})
	}

	return tags, nil
}

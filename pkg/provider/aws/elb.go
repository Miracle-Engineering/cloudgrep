package aws

import (
	"context"
	"fmt"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"

	"github.com/run-x/cloudgrep/pkg/model"
)

func (p *Provider) FetchLoadBalancers(ctx context.Context, output chan<- model.Resource) error {
	resourceType := "elb.LoadBalancer"
	elbClient := elbv2.NewFromConfig(p.config)
	input := &elbv2.DescribeLoadBalancersInput{}
	paginator := elbv2.NewDescribeLoadBalancersPaginator(elbClient, input)

	resourceConverter := p.converterFor(resourceType)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 Instances: %w", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.LoadBalancers, p.FetchLoadBalancerTag); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) FetchLoadBalancerTag(ctx context.Context, lb types.LoadBalancer) (model.Tags, error) {
	elbClient := elbv2.NewFromConfig(p.config)
	tagsResponse, err := elbClient.DescribeTags(
		ctx,
		&elbv2.DescribeTagsInput{ResourceArns: []string{*lb.LoadBalancerArn}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags for load balancer %v: %w", *lb.LoadBalancerArn, err)
	}
	var tags model.Tags
	for _, tagDescription := range tagsResponse.TagDescriptions {
		for _, tag := range tagDescription.Tags {
			tags = append(tags, model.Tag{Key: *tag.Key, Value: *tag.Value})
		}
	}
	return tags, nil
}

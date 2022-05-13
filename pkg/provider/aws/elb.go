package aws

import (
	"context"
	"fmt"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/run-x/cloudgrep/pkg/model"
)

func (p *AWSProvider) FetchLoadBalancers(ctx context.Context) ([]types.LoadBalancer, error) {
	lbOutput, err := p.elbClient.DescribeLoadBalancers(ctx, &elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}
	return lbOutput.LoadBalancers, nil
}

func (p *AWSProvider) FetchLoadBalancerTag(ctx context.Context, lb types.LoadBalancer) (model.Tags, error) {
	tagsResponse, err := p.elbClient.DescribeTags(
		ctx,
		&elbv2.DescribeTagsInput{ResourceArns: []string{*lb.LoadBalancerArn}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags for load balancer %v: %w", &lb.LoadBalancerArn, err)
	}
	var tags model.Tags
	for _, tagDescription := range tagsResponse.TagDescriptions {
		for _, tag := range tagDescription.Tags {
			tags = append(tags, model.Tag{Key: *tag.Key, Value: *tag.Value})
		}
	}
	return tags, nil
}

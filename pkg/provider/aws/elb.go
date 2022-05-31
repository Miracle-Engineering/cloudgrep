package aws

import (
	"context"
	"fmt"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
)

func (p *AWSProvider) FetchLoadBalancers(ctx context.Context, output chan<- types.LoadBalancer) error {
	pages := elbv2.NewDescribeLoadBalancersPaginator(p.elbClient, &elbv2.DescribeLoadBalancersInput{})
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch Load Balancers: %w", err)
		}

		if err := util.SendAllFromSlice(ctx, output, page.LoadBalancers); err != nil {
			return err
		}
	}

	return nil
}

func (p *AWSProvider) FetchLoadBalancerTag(ctx context.Context, lb types.LoadBalancer) (model.Tags, error) {
	tagsResponse, err := p.elbClient.DescribeTags(
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

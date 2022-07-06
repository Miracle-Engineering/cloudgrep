package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_route53(mapping map[string]mapper) {
	mapping["route53.HealthCheck"] = mapper{
		ServiceEndpointID: "route53",
		FetchFunc:         p.fetch_route53_HealthCheck,
		IdField:           "Id",
		IsGlobal:          true,
	}
	mapping["route53.HostedZone"] = mapper{
		ServiceEndpointID: "route53",
		FetchFunc:         p.fetch_route53_HostedZone,
		IdField:           "Id",
		IsGlobal:          true,
	}
}

func (p *Provider) fetch_route53_HealthCheck(ctx context.Context, output chan<- model.Resource) error {
	client := route53.NewFromConfig(p.config)
	input := &route53.ListHealthChecksInput{}

	commonTransformers := p.baseTransformers("route53.HealthCheck")
	converter := p.converterFor("route53.HealthCheck")
	transformers := append(
		resourceconverter.AllToGeneric[types.HealthCheck](commonTransformers...),
		resourceconverter.WithConverter[types.HealthCheck](converter),
		resourceconverter.WithTagFunc(p.getTags_route53_HealthCheck),
	)
	paginator := route53.NewListHealthChecksPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "route53.HealthCheck", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.HealthChecks, transformers...); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTags_route53_HealthCheck(ctx context.Context, resource types.HealthCheck) (model.Tags, error) {
	client := route53.NewFromConfig(p.config)
	input := &route53.ListTagsForResourcesInput{}

	input.ResourceIds = []string{*resource.Id}
	{
		var err error
		if err = listHealthCheckTagsInput(input); err != nil {
			return nil, fmt.Errorf("error overriding input with %s(input) for %s", "listHealthCheckTagsInput", "route53.HealthCheck")
		}
	}

	output, err := client.ListTagsForResources(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "route53.HealthCheck", err)
	}
	tagField_0 := output.ResourceTagSets
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

func (p *Provider) fetch_route53_HostedZone(ctx context.Context, output chan<- model.Resource) error {
	client := route53.NewFromConfig(p.config)
	input := &route53.ListHostedZonesInput{}

	commonTransformers := p.baseTransformers("route53.HostedZone")
	converter := p.converterFor("route53.HostedZone")
	transformers := append(
		resourceconverter.AllToGeneric[types.HostedZone](commonTransformers...),
		resourceconverter.WithConverter[types.HostedZone](converter),
		resourceconverter.WithTagFunc(p.getTags_route53_HostedZone),
	)
	paginator := route53.NewListHostedZonesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "route53.HostedZone", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.HostedZones, transformers...); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTags_route53_HostedZone(ctx context.Context, resource types.HostedZone) (model.Tags, error) {
	client := route53.NewFromConfig(p.config)
	input := &route53.ListTagsForResourcesInput{}

	input.ResourceIds = []string{*resource.Id}
	{
		var err error
		if err = listHostedZoneTagsInput(input); err != nil {
			return nil, fmt.Errorf("error overriding input with %s(input) for %s", "listHostedZoneTagsInput", "route53.HostedZone")
		}
	}

	output, err := client.ListTagsForResources(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "route53.HostedZone", err)
	}
	tagField_0 := output.ResourceTagSets
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

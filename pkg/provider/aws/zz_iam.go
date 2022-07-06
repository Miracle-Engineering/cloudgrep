package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_iam(mapping map[string]mapper) {
	mapping["iam.OpenIDConnectProvider"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_OpenIDConnectProvider,
		IdField:           "Arn",
		IsGlobal:          true,
	}
	mapping["iam.Policy"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_Policy,
		IdField:           "Arn",
		IsGlobal:          true,
	}
	mapping["iam.SAMLProvider"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_SAMLProvider,
		IdField:           "Arn",
		IsGlobal:          true,
	}
	mapping["iam.VirtualMFADevice"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_VirtualMFADevice,
		IdField:           "SerialNumber",
		IsGlobal:          true,
	}
}

func (p *Provider) fetch_iam_OpenIDConnectProvider(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListOpenIDConnectProvidersInput{}

	commonTransformers := p.baseTransformers("iam.OpenIDConnectProvider")
	converter := p.converterFor("iam.OpenIDConnectProvider")
	transformers := append(
		resourceconverter.AllToGeneric[types.OpenIDConnectProviderListEntry](commonTransformers...),
		resourceconverter.WithConverter[types.OpenIDConnectProviderListEntry](converter),
		resourceconverter.WithTagFunc(p.getTags_iam_OpenIDConnectProvider),
	)
	results, err := client.ListOpenIDConnectProviders(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "iam.OpenIDConnectProvider", err)
	}
	if err := resourceconverter.SendAll(ctx, output, results.OpenIDConnectProviderList, transformers...); err != nil {
		return err
	}

	return nil
}
func (p *Provider) getTags_iam_OpenIDConnectProvider(ctx context.Context, resource types.OpenIDConnectProviderListEntry) (model.Tags, error) {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListOpenIDConnectProviderTagsInput{}

	input.OpenIDConnectProviderArn = resource.Arn

	output, err := client.ListOpenIDConnectProviderTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "iam.OpenIDConnectProvider", err)
	}
	tagField_0 := output.Tags

	var tags model.Tags

	for _, field := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   *field.Key,
			Value: *field.Value,
		})
	}

	return tags, nil
}

func (p *Provider) fetch_iam_Policy(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListPoliciesInput{}
	input.Scope = listPoliciesScope()

	commonTransformers := p.baseTransformers("iam.Policy")
	converter := p.converterFor("iam.Policy")
	transformers := append(
		resourceconverter.AllToGeneric[types.Policy](commonTransformers...),
		resourceconverter.WithConverter[types.Policy](converter),
		resourceconverter.WithTagFunc(p.getTags_iam_Policy),
	)
	paginator := iam.NewListPoliciesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.Policy", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.Policies, transformers...); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTags_iam_Policy(ctx context.Context, resource types.Policy) (model.Tags, error) {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListPolicyTagsInput{}

	input.PolicyArn = resource.Arn

	output, err := client.ListPolicyTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "iam.Policy", err)
	}
	tagField_0 := output.Tags

	var tags model.Tags

	for _, field := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   *field.Key,
			Value: *field.Value,
		})
	}

	return tags, nil
}

func (p *Provider) fetch_iam_SAMLProvider(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListSAMLProvidersInput{}

	commonTransformers := p.baseTransformers("iam.SAMLProvider")
	converter := p.converterFor("iam.SAMLProvider")
	transformers := append(
		resourceconverter.AllToGeneric[types.SAMLProviderListEntry](commonTransformers...),
		resourceconverter.WithConverter[types.SAMLProviderListEntry](converter),
		resourceconverter.WithTagFunc(p.getTags_iam_SAMLProvider),
	)
	results, err := client.ListSAMLProviders(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "iam.SAMLProvider", err)
	}
	if err := resourceconverter.SendAll(ctx, output, results.SAMLProviderList, transformers...); err != nil {
		return err
	}

	return nil
}
func (p *Provider) getTags_iam_SAMLProvider(ctx context.Context, resource types.SAMLProviderListEntry) (model.Tags, error) {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListSAMLProviderTagsInput{}

	input.SAMLProviderArn = resource.Arn

	output, err := client.ListSAMLProviderTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "iam.SAMLProvider", err)
	}
	tagField_0 := output.Tags

	var tags model.Tags

	for _, field := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   *field.Key,
			Value: *field.Value,
		})
	}

	return tags, nil
}

func (p *Provider) fetch_iam_VirtualMFADevice(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListVirtualMFADevicesInput{}

	commonTransformers := p.baseTransformers("iam.VirtualMFADevice")
	converter := p.converterFor("iam.VirtualMFADevice")
	transformers := append(
		resourceconverter.AllToGeneric[types.VirtualMFADevice](commonTransformers...),
		resourceconverter.WithConverter[types.VirtualMFADevice](converter),
		resourceconverter.WithTagFunc(p.getTags_iam_VirtualMFADevice),
	)
	paginator := iam.NewListVirtualMFADevicesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.VirtualMFADevice", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.VirtualMFADevices, transformers...); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTags_iam_VirtualMFADevice(ctx context.Context, resource types.VirtualMFADevice) (model.Tags, error) {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListMFADeviceTagsInput{}

	input.SerialNumber = resource.SerialNumber

	output, err := client.ListMFADeviceTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "iam.VirtualMFADevice", err)
	}
	tagField_0 := output.Tags

	var tags model.Tags

	for _, field := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   *field.Key,
			Value: *field.Value,
		})
	}

	return tags, nil
}

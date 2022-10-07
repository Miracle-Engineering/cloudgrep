package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) registerIam(mapping map[string]mapper) {
	mapping["iam.OpenIDConnectProvider"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetchIamOpenIDConnectProvider,
		IdField:           "Arn",
		IsGlobal:          true,
	}
	mapping["iam.Policy"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetchIamPolicy,
		IdField:           "Arn",
		DisplayIDField:    "PolicyName",
		IsGlobal:          true,
	}
	mapping["iam.SAMLProvider"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetchIamSAMLProvider,
		IdField:           "Arn",
		IsGlobal:          true,
	}
	mapping["iam.VirtualMFADevice"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetchIamVirtualMFADevice,
		IdField:           "SerialNumber",
		IsGlobal:          true,
	}
}

func (p *Provider) fetchIamOpenIDConnectProvider(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListOpenIDConnectProvidersInput{}

	resourceConverter := p.converterFor("iam.OpenIDConnectProvider")
	var transformers resourceconverter.Transformers[types.OpenIDConnectProviderListEntry]
	transformers.AddNamed("tags", resourceconverter.TagTransformer(p.getTagsIamOpenIDConnectProvider))
	transformers.AddResource(displayIdArnPrefix("oidc-provider/"))
	results, err := client.ListOpenIDConnectProviders(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "iam.OpenIDConnectProvider", err)
	}
	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, results.OpenIDConnectProviderList, transformers); err != nil {
		return err
	}

	return nil
}
func (p *Provider) getTagsIamOpenIDConnectProvider(ctx context.Context, resource types.OpenIDConnectProviderListEntry) (model.Tags, error) {
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

func (p *Provider) fetchIamPolicy(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListPoliciesInput{}
	input.Scope = listPoliciesScope()

	resourceConverter := p.converterFor("iam.Policy")
	var transformers resourceconverter.Transformers[types.Policy]
	transformers.AddNamed("tags", resourceconverter.TagTransformer(p.getTagsIamPolicy))
	paginator := iam.NewListPoliciesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.Policy", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Policies, transformers); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTagsIamPolicy(ctx context.Context, resource types.Policy) (model.Tags, error) {
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

func (p *Provider) fetchIamSAMLProvider(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListSAMLProvidersInput{}

	resourceConverter := p.converterFor("iam.SAMLProvider")
	var transformers resourceconverter.Transformers[types.SAMLProviderListEntry]
	transformers.AddNamed("tags", resourceconverter.TagTransformer(p.getTagsIamSAMLProvider))
	transformers.AddResource(displayIdArnPrefix("saml-provider/"))
	results, err := client.ListSAMLProviders(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "iam.SAMLProvider", err)
	}
	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, results.SAMLProviderList, transformers); err != nil {
		return err
	}

	return nil
}
func (p *Provider) getTagsIamSAMLProvider(ctx context.Context, resource types.SAMLProviderListEntry) (model.Tags, error) {
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

func (p *Provider) fetchIamVirtualMFADevice(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListVirtualMFADevicesInput{}

	resourceConverter := p.converterFor("iam.VirtualMFADevice")
	var transformers resourceconverter.Transformers[types.VirtualMFADevice]
	transformers.AddNamed("tags", resourceconverter.TagTransformer(p.getTagsIamVirtualMFADevice))
	transformers.AddResource(displayIdArnPrefix("mfa/"))
	paginator := iam.NewListVirtualMFADevicesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.VirtualMFADevice", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.VirtualMFADevices, transformers); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTagsIamVirtualMFADevice(ctx context.Context, resource types.VirtualMFADevice) (model.Tags, error) {
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

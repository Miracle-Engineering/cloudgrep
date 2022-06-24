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
	mapping["iam.InstanceProfile"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_InstanceProfile,
		IdField:           "InstanceProfileId",
		IsGlobal:          true,
	}
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
	mapping["iam.Role"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_Role,
		IdField:           "RoleId",
		IsGlobal:          true,
	}
	mapping["iam.SAMLProvider"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_SAMLProvider,
		IdField:           "Arn",
		IsGlobal:          true,
	}
	mapping["iam.User"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_User,
		IdField:           "UserId",
		IsGlobal:          true,
	}
	mapping["iam.VirtualMFADevice"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_VirtualMFADevice,
		IdField:           "SerialNumber",
		IsGlobal:          true,
	}
}

func (p *Provider) fetch_iam_InstanceProfile(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListInstanceProfilesInput{}

	resourceConverter := p.converterFor("iam.InstanceProfile")
	paginator := iam.NewListInstanceProfilesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.InstanceProfile", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.InstanceProfiles, p.getTags_iam_InstanceProfile); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTags_iam_InstanceProfile(ctx context.Context, resource types.InstanceProfile) (model.Tags, error) {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListInstanceProfileTagsInput{}

	input.InstanceProfileName = resource.InstanceProfileName

	output, err := client.ListInstanceProfileTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "iam.InstanceProfile", err)
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

func (p *Provider) fetch_iam_OpenIDConnectProvider(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListOpenIDConnectProvidersInput{}

	resourceConverter := p.converterFor("iam.OpenIDConnectProvider")
	results, err := client.ListOpenIDConnectProviders(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "iam.OpenIDConnectProvider", err)
	}
	if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, results.OpenIDConnectProviderList, p.getTags_iam_OpenIDConnectProvider); err != nil {
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

	resourceConverter := p.converterFor("iam.Policy")
	paginator := iam.NewListPoliciesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.Policy", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.Policies, p.getTags_iam_Policy); err != nil {
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

func (p *Provider) fetch_iam_Role(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListRolesInput{}

	resourceConverter := p.converterFor("iam.Role")
	paginator := iam.NewListRolesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.Role", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.Roles, p.getTags_iam_Role); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTags_iam_Role(ctx context.Context, resource types.Role) (model.Tags, error) {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListRoleTagsInput{}

	input.RoleName = resource.RoleName

	output, err := client.ListRoleTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "iam.Role", err)
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

	resourceConverter := p.converterFor("iam.SAMLProvider")
	results, err := client.ListSAMLProviders(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "iam.SAMLProvider", err)
	}
	if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, results.SAMLProviderList, p.getTags_iam_SAMLProvider); err != nil {
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

func (p *Provider) fetch_iam_User(ctx context.Context, output chan<- model.Resource) error {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListUsersInput{}

	resourceConverter := p.converterFor("iam.User")
	paginator := iam.NewListUsersPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.User", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.Users, p.getTags_iam_User); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTags_iam_User(ctx context.Context, resource types.User) (model.Tags, error) {
	client := iam.NewFromConfig(p.config)
	input := &iam.ListUserTagsInput{}

	input.UserName = resource.UserName

	output, err := client.ListUserTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "iam.User", err)
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

	resourceConverter := p.converterFor("iam.VirtualMFADevice")
	paginator := iam.NewListVirtualMFADevicesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "iam.VirtualMFADevice", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.VirtualMFADevices, p.getTags_iam_VirtualMFADevice); err != nil {
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

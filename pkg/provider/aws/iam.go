package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_iam_manual(mapping map[string]mapper) {
	mapping["iam.User"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_User,
		IdField:           "UserId",
		IsGlobal:          true,
	}
	mapping["iam.InstanceProfile"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_InstanceProfile,
		IdField:           "InstanceProfileId",
		IsGlobal:          true,
	}
	mapping["iam.Role"] = mapper{
		ServiceEndpointID: "iam",
		FetchFunc:         p.fetch_iam_Role,
		IdField:           "RoleId",
		IsGlobal:          true,
	}
}

func listPoliciesScope() types.PolicyScopeType {
	return types.PolicyScopeTypeLocal
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

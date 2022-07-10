package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) registerSns(mapping map[string]mapper) {
	mapping["sns.Topic"] = mapper{
		ServiceEndpointID: "sns",
		FetchFunc:         p.fetchSnsTopic,
		IdField:           "TopicArn",
		IsGlobal:          false,
	}
}

func (p *Provider) fetchSnsTopic(ctx context.Context, output chan<- model.Resource) error {
	client := sns.NewFromConfig(p.config)
	input := &sns.ListTopicsInput{}

	resourceConverter := p.converterFor("sns.Topic")
	var transformers resourceconverter.Transformers[types.Topic]
	transformers.AddNamed("tags", resourceconverter.TagTransformer(p.getTagsSnsTopic))
	paginator := sns.NewListTopicsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "sns.Topic", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Topics, transformers); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTagsSnsTopic(ctx context.Context, resource types.Topic) (model.Tags, error) {
	client := sns.NewFromConfig(p.config)
	input := &sns.ListTagsForResourceInput{}

	input.ResourceArn = resource.TopicArn

	output, err := client.ListTagsForResource(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "sns.Topic", err)
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

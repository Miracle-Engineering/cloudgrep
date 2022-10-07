package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_cloudfront(mapping map[string]mapper) {
	mapping["cloudfront.Distribution"] = mapper{
		FetchFunc: p.fetch_cloudfront_Distribution,
		IdField:   "Id",
		IsGlobal:  true,
	}
}

func (p *Provider) fetch_cloudfront_Distribution(ctx context.Context, output chan<- model.Resource) error {
	client := cloudfront.NewFromConfig(p.config)
	resourceConverter := p.converterFor("cloudfront.Distribution")
	input := &cloudfront.ListDistributionsInput{}
	paginator := cloudfront.NewListDistributionsPaginator(client, input)
	var distributions []types.DistributionSummary

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "cloudfront.Distribution", err)
		}
		distributions = append(distributions, page.DistributionList.Items...)
	}

	var transformers resourceconverter.Transformers[types.DistributionSummary]
	transformers.AddTags(p.getTags_cloudfront_Distribution)

	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, distributions, transformers); err != nil {
		return err
	}

	return nil
}

func (p *Provider) getTags_cloudfront_Distribution(ctx context.Context, resource types.DistributionSummary) (model.Tags, error) {
	client := cloudfront.NewFromConfig(p.config)
	input := &cloudfront.ListTagsForResourceInput{Resource: resource.ARN}
	output, err := client.ListTagsForResource(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags for %s: %w", "cloudfront.Distribution", err)
	}
	var tags model.Tags

	for _, tag := range output.Tags.Items {
		tags = append(tags, model.Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}

	return tags, nil
}

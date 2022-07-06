package aws

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

const SqsQueueIdentifier = "QueueUrl"

func (p *Provider) register_sqs(mapping map[string]mapper) {
	mapping["sqs.Queue"] = mapper{
		FetchFunc: p.fetchSqsQueue,
	}
}

func (p *Provider) fetchSqsQueue(ctx context.Context, output chan<- model.Resource) error {
	const typ = "sqs.Queue"

	client := sqs.NewFromConfig(p.config)
	input := &sqs.ListQueuesInput{}

	commonTransformers := p.baseTransformers(typ)

	transformers := append(
		resourceconverter.AllToGeneric[string](commonTransformers...),
		resourceconverter.WithIDFunc(p.sqsQueueId),
		resourceconverter.WithRawDataFunc(p.sqsQueueAttributes),
		resourceconverter.WithTagFunc(p.sqsQueueTags),
	)

	paginator := sqs.NewListQueuesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "sqs.Queue", err)
		}

		err = resourceconverter.SendAll(ctx, output, page.QueueUrls, transformers...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) sqsQueueAttributes(ctx context.Context, url string) (any, error) {
	client := sqs.NewFromConfig(p.config)

	input := sqs.GetQueueAttributesInput{
		QueueUrl:       &url,
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
	}
	result, err := client.GetQueueAttributes(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", "sqs.Queue", err)
	}
	result.Attributes[SqsQueueIdentifier] = url

	return result.Attributes, nil
}

func (p *Provider) sqsQueueTags(ctx context.Context, url string) (model.Tags, error) {
	client := sqs.NewFromConfig(p.config)
	input := &sqs.ListQueueTagsInput{
		QueueUrl: &url,
	}

	output, err := client.ListQueueTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "sqs.Queue", err)
	}

	var tags model.Tags

	for key, value := range output.Tags {
		tags = append(tags, model.Tag{
			Key:   key,
			Value: value,
		})
	}

	return tags, nil
}

func (p *Provider) sqsQueueId(ctx context.Context, queueUrl string) (string, error) {
	parsed, err := url.Parse(queueUrl)
	if err != nil {
		return "", fmt.Errorf("unable to parse SQS queue URL (%s): %w", queueUrl, err)
	}

	pathParts := strings.Split(parsed.Path, "/")
	idx := len(pathParts) - 1
	if idx > 0 && pathParts[idx] == "" {
		idx--
	}

	if pathParts[idx] == "" {
		return "", fmt.Errorf("unexpected empty URL part (%s) at index %d", queueUrl, idx)
	}

	return pathParts[idx], nil
}

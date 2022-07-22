package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

const SqsQueueIdentifier = "QueueUrl"

func (p *Provider) register_sqs(mapping map[string]mapper) {
	mapping["sqs.Queue"] = mapper{
		FetchFunc:       p.fetch_sqs_Queue,
		IdField:         "QueueUrl",
		DisplayIDField:  "QueueArn",
		IsGlobal:        false,
		UseMapConverter: true,
	}
}

func (p *Provider) get_sqs_queue_urls(ctx context.Context) ([]string, error) {
	client := sqs.NewFromConfig(p.config)
	input := &sqs.ListQueuesInput{}

	paginator := sqs.NewListQueuesPaginator(client, input)
	var queueUrls []string

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s: %w", "sqs.Queue", err)
		}
		queueUrls = append(queueUrls, page.QueueUrls...)
	}
	return queueUrls, nil
}

func (p *Provider) fetch_sqs_Queue(ctx context.Context, output chan<- model.Resource) error {
	var err error
	client := sqs.NewFromConfig(p.config)
	resourceConverter := p.converterFor("sqs.Queue")
	queueUrls, err := p.get_sqs_queue_urls(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "sqs.Queue", err)
	}
	var queuesAttributes []map[string]string
	for _, queueUrl := range queueUrls {
		getQueueAttributesInput := sqs.GetQueueAttributesInput{
			QueueUrl:       &queueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
		}
		getQueueAttributesResult, err := client.GetQueueAttributes(ctx, &getQueueAttributesInput)
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "sqs.Queue", err)
		}
		getQueueAttributesResult.Attributes[SqsQueueIdentifier] = queueUrl
		queuesAttributes = append(queuesAttributes, getQueueAttributesResult.Attributes)
	}

	var transformers resourceconverter.Transformers[map[string]string]
	transformers.AddNamed("tags", resourceconverter.TagTransformer(p.getTags_sqs_Queue))
	transformers.AddNamedResource("displayId", displayIdArn)

	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, queuesAttributes, transformers); err != nil {
		return err
	}
	return nil
}

func (p *Provider) getTags_sqs_Queue(ctx context.Context, resource map[string]string) (model.Tags, error) {
	queueUrl := resource[SqsQueueIdentifier]
	client := sqs.NewFromConfig(p.config)
	input := &sqs.ListQueueTagsInput{
		QueueUrl: &queueUrl,
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

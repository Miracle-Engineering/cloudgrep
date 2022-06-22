package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/hashicorp/go-multierror"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

const SQS_QUEUE_IDENTIFIER = "QueueUrl"

func (p *Provider) register_sqs(mapping map[string]mapper) {
	mapping["sqs.SQS"] = mapper{
		FetchFunc:       p.fetch_sqs_SQS,
		IdField:         "QueueUrl",
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
			return nil, fmt.Errorf("failed to fetch %s: %w", "sqs.SQS", err)
		}
		queueUrls = append(queueUrls, page.QueueUrls...)
	}
	return queueUrls, nil
}

func (p *Provider) fetch_sqs_SQS(ctx context.Context, output chan<- model.Resource) error {
	var err error
	client := sqs.NewFromConfig(p.config)
	resourceConverter := p.converterFor("sqs.SQS")
	queueUrls, err := p.get_sqs_queue_urls(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "sqs.SQS", err)
	}
	var errors *multierror.Error
	var queuesAttributes []map[string]string
	for _, queueUrl := range queueUrls {
		getQueueAttributesInput := sqs.GetQueueAttributesInput{
			QueueUrl:       &queueUrl,
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
		}
		getQueueAttributesResult, err := client.GetQueueAttributes(ctx, &getQueueAttributesInput)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
		getQueueAttributesResult.Attributes[SQS_QUEUE_IDENTIFIER] = queueUrl
		queuesAttributes = append(queuesAttributes, getQueueAttributesResult.Attributes)
	}
	if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, queuesAttributes, p.getTags_sqs_SQS); err != nil {
		return err
	}
	return nil
}

func (p *Provider) getTags_sqs_SQS(ctx context.Context, resource map[string]string) (model.Tags, error) {
	queueUrl := resource[SQS_QUEUE_IDENTIFIER]
	client := sqs.NewFromConfig(p.config)
	input := &sqs.ListQueueTagsInput{
		QueueUrl: &queueUrl,
	}
	output, err := client.ListQueueTags(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "sqs.SQS", err)
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

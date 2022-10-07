package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_s3(mapping map[string]mapper) {
	mapping["s3.Bucket"] = mapper{
		FetchFunc: p.fetch_s3_Bucket,
		IdField:   "Name",
		IsGlobal:  true,
	}
}

func (p *Provider) fetch_s3_Bucket(ctx context.Context, output chan<- model.Resource) error {
	client := s3.NewFromConfig(p.config)
	input := &s3.ListBucketsInput{}

	var transformers resourceconverter.Transformers[types.Bucket]
	transformers.AddTags(p.getTags_s3_Bucket)

	resourceConverter := p.converterFor("s3.Bucket")
	results, err := client.ListBuckets(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "s3.Bucket", err)
	}
	if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, results.Buckets, transformers); err != nil {
		return err
	}

	return nil
}
func (p *Provider) getTags_s3_Bucket(ctx context.Context, resource types.Bucket) (model.Tags, error) {
	client := s3.NewFromConfig(p.config)
	locationInput := &s3.GetBucketLocationInput{Bucket: resource.Name}
	locationOutput, err := client.GetBucketLocation(ctx, locationInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s location: %w", "s3.Bucket", err)
	}
	tempConfig := p.config.Copy()
	switch location := locationOutput.LocationConstraint; location {
	case "":
		tempConfig.Region = "us-east-1"
	default:
		tempConfig.Region = string(location)
	}
	client = s3.NewFromConfig(tempConfig)

	input := &s3.GetBucketTaggingInput{}
	input.Bucket = resource.Name

	output, err := client.GetBucketTagging(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NoSuchTagSet" {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "s3.Bucket", err)
	}
	tagField_0 := output.TagSet

	var tags model.Tags

	for _, field := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   *field.Key,
			Value: *field.Value,
		})
	}

	return tags, nil
}

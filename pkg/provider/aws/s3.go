package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
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

	resourceConverter := p.converterFor("s3.Bucket")
	commonTransformers := p.baseTransformers("s3.Bucket")
	transformers := append(
		resourceconverter.AllToGeneric[types.Bucket](commonTransformers...),
		resourceconverter.WithConverter[types.Bucket](resourceConverter),
		resourceconverter.WithRegionFunc(p.getS3BucketRegion),
		p.getTags_s3_Bucket,
	)

	results, err := client.ListBuckets(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "s3.Bucket", err)
	}
	if err := resourceconverter.SendAll(ctx, output, results.Buckets, transformers...); err != nil {
		return err
	}

	return nil
}

func (p *Provider) getS3BucketRegion(ctx context.Context, bucket types.Bucket) (string, error) {
	client := s3.NewFromConfig(p.config)
	input := &s3.GetBucketLocationInput{Bucket: bucket.Name}
	output, err := client.GetBucketLocation(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get %s location: %w", "s3.Bucket", err)
	}

	if output.LocationConstraint == "" {
		return "us-east-1", nil
	}

	// TODO: Should we have special handling for the "EU" location constraint?

	return string(output.LocationConstraint), nil
}

func (p *Provider) getTags_s3_Bucket(ctx context.Context, resource types.Bucket, res *model.Resource) error {
	config := p.config.Copy()
	config.Region = res.Region
	client := s3.NewFromConfig(config)

	input := &s3.GetBucketTaggingInput{}
	input.Bucket = resource.Name

	output, err := client.GetBucketTagging(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NoSuchTagSet" {
				return nil
			}
		}
		return fmt.Errorf("failed to fetch %s tags: %w", "s3.Bucket", err)
	}
	tagField_0 := output.TagSet

	var tags model.Tags

	for _, field := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   *field.Key,
			Value: *field.Value,
		})
	}

	res.Tags = append(res.Tags, tags...)

	return nil
}

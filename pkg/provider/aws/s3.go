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

func (p *Provider) FetchS3Buckets(ctx context.Context, output chan<- model.Resource) error {
	resourceType := "s3.Bucket"
	s3Client := s3.NewFromConfig(p.config)
	input := &s3.ListBucketsInput{}
	result, err := s3Client.ListBuckets(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch S3 buckets: %w", err)
	}

	resourceConverter := p.converterFor(resourceType)
	if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, result.Buckets, p.FetchS3BucketsTag); err != nil {
		return err
	}

	return nil
}

func (p *Provider) FetchS3BucketsTag(ctx context.Context, bucket types.Bucket) (model.Tags, error) {
	s3Client := s3.NewFromConfig(p.config)
	output, err := s3Client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{Bucket: bucket.Name})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchTagSet" {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("failed to fetch tags for S3 bucket %v: %w", *bucket.Name, err)
	}
	var tags model.Tags
	for _, t := range output.TagSet {
		tags = append(tags, model.Tag{Key: *t.Key, Value: *t.Value})
	}
	return tags, nil
}

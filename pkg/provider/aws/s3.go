package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
)

func (p *AWSProvider) FetchS3Buckets(ctx context.Context, output chan<- types.Bucket) error {
	input := &s3.ListBucketsInput{}
	result, err := p.s3Client.ListBuckets(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch S3 buckets: %w", err)
	}

	if err := util.SendAllFromSlice(ctx, output, result.Buckets); err != nil {
		return err
	}

	return nil
}

func (p *AWSProvider) FetchS3BucketsTag(ctx context.Context, bucket types.Bucket) (model.Tags, error) {
	output, err := p.s3Client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{Bucket: bucket.Name}, func(options *s3.Options) {
		options.Region = p.Region()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags for S3 bucket %v: %w", *bucket.Name, err)
	}
	var tags model.Tags
	for _, t := range output.TagSet {
		tags = append(tags, model.Tag{Key: *t.Key, Value: *t.Value})
	}
	return tags, nil
}

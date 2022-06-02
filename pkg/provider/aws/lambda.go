package aws

import (
	"context"
	"fmt"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/run-x/cloudgrep/pkg/model"
)

func (p *Provider) FetchLambdaFunctions(ctx context.Context, output chan<- model.Resource) error {
	resourceType := "lambda.Function"
	lambdaClient := lambda.NewFromConfig(p.config)
	input := &lambda.ListFunctionsInput{}
	paginator := lambda.NewListFunctionsPaginator(lambdaClient, input)

	resourceConverter := p.converterFor(resourceType)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch Lambda Function: %w", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.Functions, p.FetchLambdaFunctionTag); err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) FetchLambdaFunctionTag(ctx context.Context, fn types.FunctionConfiguration) (model.Tags, error) {
	lambdaClient := lambda.NewFromConfig(p.config)
	tagsResponse, err := lambdaClient.GetFunction(
		ctx,
		&lambda.GetFunctionInput{FunctionName: fn.FunctionArn},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags for lambda function %v: %w", *fn.FunctionArn, err)
	}
	var tags model.Tags
	for key, value := range tagsResponse.Tags {
		tags = append(tags, model.Tag{Key: key, Value: value})
	}
	return tags, nil
}

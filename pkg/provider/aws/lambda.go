package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/run-x/cloudgrep/pkg/model"
)

func (awsPrv *AWSProvider) FetchFunctions(ctx context.Context) ([]types.FunctionConfiguration, error) {
	input := &lambda.ListFunctionsInput{}
	var functions []types.FunctionConfiguration
	paginator := lambda.NewListFunctionsPaginator(awsPrv.lambdaClient, input)
	for paginator.HasMorePages() {
		result, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Lambda Functions: %w", err)
		}
		functions = append(functions, result.Functions...)
	}
	return functions, nil
}

func (p *AWSProvider) FetchFunctionTag(ctx context.Context, fn types.FunctionConfiguration) (model.Tags, error) {
	tagsResponse, err := p.lambdaClient.GetFunction(
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

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) registerLambda(mapping map[string]mapper) {
	mapping["lambda.Function"] = mapper{
		ServiceEndpointID: "lambda",
		FetchFunc:         p.fetchLambdaFunction,
		IdField:           "FunctionArn",
		DisplayIDField:    "FunctionName",
		IsGlobal:          false,
	}
}

func (p *Provider) fetchLambdaFunction(ctx context.Context, output chan<- model.Resource) error {
	client := lambda.NewFromConfig(p.config)
	input := &lambda.ListFunctionsInput{}

	resourceConverter := p.converterFor("lambda.Function")
	var transformers resourceconverter.Transformers[types.FunctionConfiguration]
	transformers.AddNamed("tags", resourceconverter.TagTransformer(p.getTagsLambdaFunction))
	paginator := lambda.NewListFunctionsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "lambda.Function", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Functions, transformers); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTagsLambdaFunction(ctx context.Context, resource types.FunctionConfiguration) (model.Tags, error) {
	client := lambda.NewFromConfig(p.config)
	input := &lambda.GetFunctionInput{}

	input.FunctionName = resource.FunctionArn

	output, err := client.GetFunction(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "lambda.Function", err)
	}
	tagField_0 := output.Tags

	var tags model.Tags

	for key, value := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   key,
			Value: value,
		})
	}

	return tags, nil
}

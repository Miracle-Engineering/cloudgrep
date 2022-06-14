package aws

import (
	"context"
	"fmt"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

type converterForer interface {
	converterFor(string) resourceconverter.ResourceConverter
}

type Paginator[Options any, Output any] interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*Options)) (*Output, error)
}

type PageOutputConverter[Output any, SDKType any] func(*Output) []SDKType

func sendPages[Options any, Output any, SDKType any](
	ctx context.Context,
	name string,
	converters converterForer,
	paginator Paginator[Options, Output],
	conv PageOutputConverter[Output, SDKType],
	output chan<- model.Resource,
	tagFunc func(context.Context, SDKType) (model.Tags, error),
) error {
	resourceConverter := converters.converterFor("elb.LoadBalancer")
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", name, err)
		}

		resources := conv(page)

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, resources, tagFunc); err != nil {
			return err
		}
	}

	return nil
}

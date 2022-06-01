package aws

import (
	"context"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
	"github.com/run-x/cloudgrep/pkg/util"
)

func SendAllConverted[T any](ctx context.Context, output chan<- model.Resource, converter resourceconverter.ResourceConverter, resources []T) error {
	var converted []model.Resource

	for _, raw := range resources {
		resource, err := converter.ToResource(ctx, raw, nil)
		if err != nil {
			return err
		}

		converted = append(converted, resource)
	}

	return util.SendAllFromSlice(ctx, output, converted)
}

type tagFunc[T any] func(context.Context, T) (model.Tags, error)

func SendAllConvertedTags[T any](ctx context.Context, output chan<- model.Resource, converter resourceconverter.ResourceConverter, resources []T, tF tagFunc[T]) error {
	var converted []model.Resource

	for _, raw := range resources {
		tags, err := tF(ctx, raw)
		if err != nil {
			return err
		}
		resource, err := converter.ToResource(ctx, raw, tags)
		if err != nil {
			return err
		}

		converted = append(converted, resource)
	}

	return util.SendAllFromSlice(ctx, output, converted)
}

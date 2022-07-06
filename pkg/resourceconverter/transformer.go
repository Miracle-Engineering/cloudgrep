package resourceconverter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/run-x/cloudgrep/pkg/model"
)

func Transform[T any](ctx context.Context, raw T, resource *model.Resource, transformers ...GenericTransformer[T]) error {
	for _, transformer := range transformers {
		err := transformer(ctx, raw, resource)
		if err != nil {
			if errors.Is(err, ctx.Err()) {
				return err
			}

			var zero T
			t := reflect.TypeOf(zero)

			return fmt.Errorf("unable to apply transformer for %s.%s (type %v, id %v): %w", t.PkgPath(), t.Name(), resource.Type, resource.Id, err)
		}
	}

	return nil
}

type ResourceTransformer func(context.Context, *model.Resource) error
type GenericTransformer[T any] func(context.Context, T, *model.Resource) error

func WithConverter[T any](rc ResourceConverter) GenericTransformer[T] {
	return func(ctx context.Context, raw T, res *model.Resource) error {
		return rc.ToResource(ctx, raw, res)
	}
}

func WithConverterDetails[T any, U any](rc ResourceConverter, detailFunc func(context.Context, T) (U, error)) GenericTransformer[T] {
	return func(ctx context.Context, raw T, res *model.Resource) error {
		details, err := detailFunc(ctx, raw)
		if err != nil {
			return err
		}

		return rc.ToResource(ctx, details, res)
	}
}

func WithTagFunc[T any](tagFunc func(context.Context, T) (model.Tags, error)) GenericTransformer[T] {
	return func(ctx context.Context, raw T, res *model.Resource) error {
		tags, err := tagFunc(ctx, raw)
		if err != nil {
			return err
		}

		res.Tags = append(res.Tags, tags...)

		return nil
	}
}

func WithIDFunc[T any](idFunc func(context.Context, T) (string, error)) GenericTransformer[T] {
	return func(ctx context.Context, raw T, res *model.Resource) error {
		id, err := idFunc(ctx, raw)
		if err != nil {
			return err
		}

		res.Id = id

		return nil
	}
}

func WithRawDataFunc[T any](rawFunc func(context.Context, T) (any, error)) GenericTransformer[T] {
	return func(ctx context.Context, raw T, res *model.Resource) error {
		data, err := rawFunc(ctx, raw)
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(data)
		if err != nil {
			return err
		}

		res.RawData = encoded

		return nil
	}
}

func WithIDTransformer(idFunc func(string) (string, error)) ResourceTransformer {
	return func(ctx context.Context, res *model.Resource) error {
		id, err := idFunc(res.Id)
		if err != nil {
			return err
		}

		res.Id = id

		return nil
	}
}

func WithRegion(region string) ResourceTransformer {
	return func(ctx context.Context, res *model.Resource) error {
		res.Region = region

		return nil
	}
}

func WithType(typ string) ResourceTransformer {
	return func(ctx context.Context, res *model.Resource) error {
		res.Type = typ

		return nil
	}
}

func ToGeneric[T any](t ResourceTransformer) GenericTransformer[T] {
	return func(ctx context.Context, _ T, res *model.Resource) error {
		return t(ctx, res)
	}
}

func AllToGeneric[T any](in ...ResourceTransformer) []GenericTransformer[T] {
	var converted []GenericTransformer[T]
	for _, t := range in {
		converted = append(converted, ToGeneric[T](t))
	}

	return converted
}

func SendAll[T any](ctx context.Context, output chan<- model.Resource, resources []T, transformers ...GenericTransformer[T]) error {
	for _, raw := range resources {
		resource := model.Resource{}

		if err := Transform(ctx, raw, &resource, transformers...); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- resource:
		}
	}

	return nil
}

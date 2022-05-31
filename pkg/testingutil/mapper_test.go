package testingutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
	"github.com/stretchr/testify/assert"
)

type MockInstance struct {
	Wrapped model.Resource
}

func fetchTestInstances(ctx context.Context, output chan<- MockInstance) error {
	resources := []MockInstance{
		{
			Wrapped: model.Resource{
				Id: "i-1",
				Tags: []model.Tag{
					{Key: "IntegrationTest", Value: "true"},
					{Key: "tag2", Value: "val2"},
				},
			},
		},
		{
			Wrapped: model.Resource{
				Id: "i-2",
				Tags: []model.Tag{
					{Key: "tag2", Value: "val2"},
				},
			},
		},
	}

	return util.SendAllFromSlice(ctx, output, resources)
}

func TestConvertToResources(t *testing.T) {
	raw := MustFetchAll(context.Background(), t, fetchTestInstances)

	mapper := &FakeMapper{}

	fake := Fake(t)

	resources := ConvertToResources(fake, context.Background(), mapper, raw)
	assert.False(t, fake.IsFail)
	assert.Empty(t, fake.Logs)
	assert.Len(t, resources, 1)
}

type FakeMapper struct {
	Err error
}

func (m *FakeMapper) ToResource(_ context.Context, raw any, region string) (model.Resource, error) {
	if m.Err != nil {
		return model.Resource{}, m.Err
	}

	instance := raw.(MockInstance)

	return instance.Wrapped, nil
}

func TestConvertToResources_error(t *testing.T) {
	raw := MustFetchAll(context.Background(), t, fetchTestInstances)

	mapper := &FakeMapper{
		Err: fmt.Errorf("failed"),
	}

	fake := Fake(t)

	assert.Panics(t, func() {
		ConvertToResources(fake, context.Background(), mapper, raw)
	})
}

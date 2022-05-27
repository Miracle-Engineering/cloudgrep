package testingutil

import (
	"context"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
	"github.com/stretchr/testify/assert"
)

const TestRegion = "us-east-1"

func AssertResourceCount(t testing.TB, resources []model.Resource, tagValue string, count int) {
	t.Helper()
	if tagValue == "" {
		resources = ResourceFilterTagKey(resources, "test")
	} else {
		resources = ResourceFilterTagKeyValue(resources, "test", tagValue)
	}

	assert.Lenf(t, resources, count, "expected %d resource(s) with tag test=%s", count, tagValue)
}

func ConvertToResources[T any](t testing.TB, ctx context.Context, mapper mapper.Mapper, raw []T) []model.Resource {
	t.Helper()
	output := make([]model.Resource, 0, len(raw))
	for _, in := range raw {
		resource, err := mapper.ToResource(ctx, in, TestRegion)
		if err != nil {
			t.Fatalf("cannot convert resource: %v", err)
		}

		output = append(output, resource)
	}

	// Make sure we only grab resources meant for integration testing
	output = ResourceFilterTagKeyValue(output, "IntegrationTest", "true")

	return output
}

func ResourceFilterTagKey(in []model.Resource, key string) []model.Resource {
	return FilterFunc(in, func(r model.Resource) bool {
		for _, tag := range r.Tags {
			if tag.Key == key {
				return true
			}
		}

		return false
	})
}

func ResourceFilterTagKeyValue(in []model.Resource, key, value string) []model.Resource {
	return FilterFunc(in, func(r model.Resource) bool {
		for _, tag := range r.Tags {
			if tag.Key == key && tag.Value == value {
				return true
			}
		}

		return false
	})
}

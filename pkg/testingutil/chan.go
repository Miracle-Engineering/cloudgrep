package testingutil

import (
	"context"
	"testing"
)

type FetchAllFunc[T any] func(context.Context, chan<- T) error

// FetchAll pulls all resources from `f` over the passed channel, returning the resources as a slice
func FetchAll[T any](ctx context.Context, t testing.TB, f FetchAllFunc[T]) ([]T, error) {
	t.Helper()

	var resources []T
	resourceChan := make(chan T)
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		for r := range resourceChan {
			resources = append(resources, r)
		}
	}()

	err := f(ctx, resourceChan)
	close(resourceChan)
	<-doneCh

	return resources, err
}

// MustFetchAll is like FetchAll, but fatals the running test if there is an error during fetching
func MustFetchAll[T any](ctx context.Context, t testing.TB, f FetchAllFunc[T]) []T {
	t.Helper()

	resources, err := FetchAll(ctx, t, f)
	if err != nil {
		t.Fatalf("error with testingutil.FetchAll: %v", err)
	}

	return resources
}

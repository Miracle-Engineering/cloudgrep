package testingutil

import (
	"context"
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestFetchAll(t *testing.T) {
	in := []string{
		"foo",
		"bar",
	}

	fetchFunc := func(ctx context.Context, out chan<- string) error {
		return util.SendAllFromSlice(ctx, out, in)
	}

	actual, err := FetchAll(context.Background(), t, fetchFunc)

	assert.NoError(t, err)
	assert.Equal(t, in, actual)
}

func TestMustFetchAll(t *testing.T) {
	in := []string{
		"foo",
		"bar",
	}

	fetchFunc := func(ctx context.Context, out chan<- string) error {
		return util.SendAllFromSlice(ctx, out, in)
	}

	actual := MustFetchAll(context.Background(), t, fetchFunc)

	assert.Equal(t, in, actual)
}

func TestFetchAllCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	in := []string{
		"foo",
		"bar",
	}

	fetchFunc := func(ctx context.Context, out chan<- string) error {
		return util.SendAllFromSlice(ctx, out, in)
	}

	actual, err := FetchAll(ctx, t, fetchFunc)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, actual)
}

func TestMustFetchAllCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	in := []string{
		"foo",
		"bar",
	}

	fetchFunc := func(ctx context.Context, out chan<- string) error {
		return util.SendAllFromSlice(ctx, out, in)
	}

	fake := Fake(t)
	MustFetchAll(ctx, fake, fetchFunc)
	assert.Contains(t, fake.Logs, "error with testingutil.FetchAll: context canceled")
	assert.True(t, fake.IsFail)
}

package datastore

import (
	"context"
	"errors"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestBlackholeStats(t *testing.T) {
	ctx := context.Background()
	resources := []*model.Resource{
		{},
		{},
		{},
	}

	ds := &Blackhole{}
	err := ds.WriteResources(ctx, resources)
	assert.NoError(t, err)

	stats, err := ds.Stats(ctx)

	assert.NoError(t, err)
	assert.Equal(t, len(resources), stats.ResourcesCount)
	assert.Equal(t, len(resources), ds.Count())
}

func TestBlackholeEmptyFuncs(t *testing.T) {
	// Make sure all funcs in the blackhole store have coverage

	ctx := context.Background()
	ds := &Blackhole{}

	var err error

	err = ds.Ping()
	assert.NoError(t, err)

	resource, err := ds.GetResource(ctx, "")
	assert.NoError(t, err)
	assert.Nil(t, resource)

	resources, err := ds.GetResources(ctx, nil)
	assert.NoError(t, err)
	assert.Empty(t, resources)

	fields, err := ds.GetFields(ctx)
	assert.NoError(t, err)
	assert.Empty(t, fields)
}

func TestBlackholeWriteError(t *testing.T) {
	ctx := context.Background()
	resources := []*model.Resource{
		{},
	}

	ds := &Blackhole{}
	err := ds.WriteResources(ctx, resources)
	assert.NoError(t, err)

	expectedErr := errors.New("foo")
	ds.SetWriteError(expectedErr)

	err = ds.WriteResources(ctx, resources)
	assert.ErrorIs(t, err, expectedErr)

	ds.SetWriteError(nil)
	err = ds.WriteResources(ctx, resources)
	assert.NoError(t, err)
}

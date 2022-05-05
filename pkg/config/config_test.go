package config

import (
	"context"
	"testing"

	"github.com/run-x/cloudgrep/pkg/options"
	"github.com/stretchr/testify/assert"
)

func Test_Init_Default(t *testing.T) {
	config, err := New(context.Background(), options.Options{})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(config.Providers))
	assert.Equal(t, "aws", config.Providers[0].Cloud)

	assert.Equal(t, "sqlite", config.Datastore.Type)
	assert.Equal(t, "file::memory:?cache=shared", config.Datastore.DataSourceName)

	assert.Equal(t, "localhost", config.Web.Host)
	assert.Equal(t, 8080, config.Web.Port)
	assert.Equal(t, "", config.Web.Prefix)

	assert.Equal(t, "dev", config.Logging.Mode)
	assert.Equal(t, true, config.Logging.IsDev())
}

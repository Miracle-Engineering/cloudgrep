package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/run-x/cloudgrep/pkg/config"
	providerutil "github.com/run-x/cloudgrep/pkg/testingutil/provider"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewProviders(t *testing.T) {

	ctx := context.Background()
	awsCfg := config.Provider{
		Cloud:   "aws",
		Regions: []string{"invalid-region"},
	}
	logger := zaptest.NewLogger(t)

	_, err := NewProviders(ctx, awsCfg, logger)
	//if it can connect, it will error with invalid AWS region
	require.True(t,
		strings.Contains(err.Error(), "invalid AWS region: invalid-region") ||
			strings.Contains(err.Error(), "no AWS credentials found"))

	fakeCfg := config.Provider{
		Cloud: "fake",
	}
	_, err = NewProviders(ctx, fakeCfg, logger)
	require.ErrorContains(t, err, "unknown provider cloud 'fake'")

	//register the fake provider
	RegisterExtraProviders("fake", fakeProviders())
	providers, err := NewProviders(ctx, fakeCfg, logger)
	require.NoError(t, err)
	require.Equal(t, "foo", providers[0].AccountId())

}

func fakeProviders() []Provider {
	fakeProviders := []*providerutil.FakeProvider{
		{
			ID: "foo",
			Foo: providerutil.FakeProviderResourceConfig{
				Count: 1,
			},
		},
	}

	providers := make([]Provider, len(fakeProviders))
	for idx, provider := range fakeProviders {
		providers[idx] = provider
	}

	return providers
}

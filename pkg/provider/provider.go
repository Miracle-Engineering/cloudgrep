package provider

import (
	"context"
	"fmt"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/provider/aws"
	"github.com/run-x/cloudgrep/pkg/provider/types"
	"go.uber.org/zap"
)

type Provider = types.Provider
type FetchFunc = types.FetchFunc

//allow dynamically registering extra providers (only used for testing at this time)
var extraProviders map[string][]Provider = map[string][]Provider{}

func NewProviders(ctx context.Context, config config.Provider, logger *zap.Logger) ([]Provider, error) {
	if config.Cloud == "aws" {
		return aws.NewProviders(ctx, config, logger)
	}
	if providers, ok := extraProviders[config.Cloud]; ok {
		return providers, nil
	}
	return nil, fmt.Errorf("unknown provider cloud '%v'", config.Cloud)
}

func RegisterExtraProviders(cloud string, providers []Provider) {
	extraProviders[cloud] = providers
}

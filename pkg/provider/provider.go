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

func NewProviders(ctx context.Context, config config.Provider, logger *zap.Logger) ([]Provider, error) {
	if config.Cloud == "aws" {
		return aws.NewProviders(ctx, config, logger)
	}
	return nil, fmt.Errorf("unknown provider cloud '%v'", config.Cloud)
}

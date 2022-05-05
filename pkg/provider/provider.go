package provider

import (
	"context"
	"fmt"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/provider/aws"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
	"go.uber.org/zap"
)

//Provider is an interface to be implemented for a cloud provider to fetch resources
//The provider must provide a mapping configuration which references the methods to fetch the resources.
//These methods need to be implemented and they will be called by a Mapper using reflection.
type Provider interface {
	GetMapper() mapper.Mapper
	Region() string
}

func NewProvider(ctx context.Context, config config.Provider, logger *zap.Logger) (Provider, error) {
	if config.Cloud == "aws" {
		return aws.NewAWSProvider(ctx, config, logger)
	}
	return nil, fmt.Errorf("unknown provider cloud '%v'", config.Cloud)
}

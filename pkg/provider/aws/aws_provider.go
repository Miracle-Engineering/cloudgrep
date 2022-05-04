package aws

import (
	"context"
	"embed"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"go.uber.org/zap"

	cfg "github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
)

type AWSProvider struct {
	logger       *zap.Logger
	config       aws.Config
	ec2Client    *ec2.Client
	mapperConfig mapper.Config
}

//go:embed mapping.yaml
var embedConfig embed.FS

func NewAWSProvider(ctx context.Context, cfg cfg.Provider, logger *zap.Logger) (*AWSProvider, error) {
	provider := AWSProvider{}
	provider.logger = logger
	//create the clients
	var err error
	logger.Info("Connecting to AWS account")
	provider.config, err = config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot load default config: %w", err)
	}
	provider.ec2Client = ec2.NewFromConfig(provider.config)
	logger.Sugar().Infow("AWS", "region", provider.Region())

	//load the mapping configuration
	data, err := embedConfig.ReadFile("mapping.yaml")
	if err != nil {
		return nil, err
	}
	provider.mapperConfig, err = mapper.LoadConfig(data)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (p AWSProvider) Region() string {
	return p.config.Region
}

func (p AWSProvider) GetMapperConfig() mapper.Config {
	return p.mapperConfig
}

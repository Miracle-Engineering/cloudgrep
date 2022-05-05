package aws

import (
	"context"
	_ "embed"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"go.uber.org/zap"

	cfg "github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
)

type AWSProvider struct {
	logger    *zap.Logger
	config    aws.Config
	ec2Client *ec2.Client
	mapper    mapper.Mapper
}

//go:embed mapping.yaml
var embedConfig []byte

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

	//create the mapper for this provider
	provider.mapper, err = mapper.New(embedConfig, *logger, reflect.ValueOf(&provider))
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (p AWSProvider) Region() string {
	return p.config.Region
}

func (p AWSProvider) GetMapper() mapper.Mapper {
	return p.mapper
}

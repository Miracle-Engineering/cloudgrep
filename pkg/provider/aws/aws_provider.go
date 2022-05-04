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
	logger    *zap.Logger
	config    aws.Config
	ec2Client *ec2.Client
}

//go:embed mapping.yaml
var embedConfig embed.FS

func (p *AWSProvider) Region() string {
	return p.config.Region
}

func (p *AWSProvider) Init(ctx context.Context, cfg cfg.Provider, logger *zap.Logger) (mapper.Config, error) {
	p.logger = logger
	//create the clients
	var err error
	logger.Info("Connecting to AWS account")
	p.config, err = config.LoadDefaultConfig(ctx)
	if err != nil {
		return mapper.Config{}, fmt.Errorf("cannot load default config: %w", err)
	}
	p.ec2Client = ec2.NewFromConfig(p.config)
	logger.Sugar().Infow("AWS", "region", p.Region())

	//load the mapping configuration
	data, err := embedConfig.ReadFile("mapping.yaml")
	if err != nil {
		return mapper.Config{}, err
	}
	return mapper.LoadConfig(data)
}

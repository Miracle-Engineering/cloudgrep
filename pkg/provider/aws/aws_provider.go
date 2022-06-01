package aws

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
	"github.com/run-x/cloudgrep/pkg/util"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"

	cfg "github.com/run-x/cloudgrep/pkg/config"
)

type AWSProvider struct {
	logger       *zap.Logger
	config       aws.Config
	ec2Client    *ec2.Client
	elbClient    *elbv2.Client
	s3Client     *s3.Client
	rdsClient    *rds.Client
	mapper       mapper.Mapper
	lambdaClient *lambda.Client
}

//go:embed mapping.yaml
var embedConfig []byte

func NewAWSProvider(ctx context.Context, cfg cfg.Provider, logger *zap.Logger) (*AWSProvider, error) {
	provider := AWSProvider{}
	provider.logger = logger
	var err error
	logger.Info("Connecting to AWS account")
	provider.config, err = config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot load default config: %w", err)
	}
	logger.Sugar().Infow("AWS", "region", provider.Region())
	stsClient := sts.NewFromConfig(provider.config)
	input := &sts.GetCallerIdentityInput{}

	result, err := stsClient.GetCallerIdentity(ctx, input)
	if err != nil {
		if serr, ok := err.(*smithy.OperationError); ok {
			return nil, util.NewUserError(
				fmt.Sprintf(
					"Encountered the following error when trying to verify AWS credentials: %v",
					serr.Unwrap().Error()))
		} else {
			return nil, err
		}
	}
	logger.Sugar().Infof("Using the following identity: %v", *result.Arn)

	//create the clients
	provider.ec2Client = ec2.NewFromConfig(provider.config)
	provider.elbClient = elbv2.NewFromConfig(provider.config)
	provider.s3Client = s3.NewFromConfig(provider.config)
	provider.lambdaClient = lambda.NewFromConfig(provider.config)
	provider.rdsClient = rds.NewFromConfig(provider.config)

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

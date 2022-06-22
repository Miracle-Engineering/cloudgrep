package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	cfg "github.com/run-x/cloudgrep/pkg/config"
	regionutil "github.com/run-x/cloudgrep/pkg/provider/aws/regions"
	"github.com/run-x/cloudgrep/pkg/provider/types"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
	"github.com/run-x/cloudgrep/pkg/util"
	_ "github.com/run-x/cloudgrep/pkg/util/rlimit"
	"go.uber.org/zap"
)

type Provider struct {
	config    aws.Config
	accountId string
	region    regionutil.Region
}

func (p Provider) String() string {
	return fmt.Sprintf("AWS Provider for account %v, region %v", p.accountId, p.region.ID())
}

func (p Provider) FetchFunctions() map[string]types.FetchFunc {
	funcMap := make(map[string]types.FetchFunc)
	for resourceType, mapping := range p.getTypeMapping() {
		if p.region.IsGlobal() != mapping.IsGlobal {
			continue
		}

		if mapping.ServiceEndpointID != "" && !p.region.IsServiceSupported(mapping.ServiceEndpointID) {
			continue
		}

		funcMap[resourceType] = mapping.FetchFunc
	}
	return funcMap
}

func (p *Provider) converterFor(resourceType string) resourceconverter.ResourceConverter {
	mapping, ok := p.getTypeMapping()[resourceType]
	if !ok {
		panic(fmt.Sprintf("Could not find mapping for resource type %v", resourceType))
	}

	region := p.region.ID()
	if mapping.UseMapConverter {
		return &resourceconverter.MapConverter{
			Region:       region,
			ResourceType: resourceType,
			TagField:     mapping.TagField,
			IdField:      mapping.IdField,
		}
	}
	return &resourceconverter.ReflectionConverter{
		Region:       region,
		ResourceType: resourceType,
		TagField:     mapping.TagField,
		IdField:      mapping.IdField,
	}
}

func NewProviders(ctx context.Context, cfg cfg.Provider, logger *zap.Logger) ([]types.Provider, error) {
	logger.Info("Connecting to AWS account")
	defaultConfig, err := config.LoadDefaultConfig(ctx, config.WithDefaultsMode(aws.DefaultsModeCrossRegion))
	if err != nil {
		return nil, err
	}

	regions, err := regionutil.SelectRegions(ctx, cfg.Regions, defaultConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot select regions for AWS provider: %w", err)
	}

	regionutil.SetConfigRegion(&defaultConfig, regions)

	identity, err := VerifyCreds(ctx, defaultConfig)
	if err != nil {
		return nil, err
	}
	logger.Sugar().Infof("Using the following identity: %v", *identity.Arn)
	logger.Sugar().Infof("Will look in regions %v", regions)
	var providers []types.Provider

	for _, region := range regions {
		newConfig := defaultConfig.Copy()
		if !region.IsGlobal() {
			newConfig.Region = region.ID()
		}

		logger.Sugar().Infof("Creating provider for AWS region %v", region)
		newProvider := Provider{
			config:    newConfig,
			accountId: *identity.Account,
			region:    region,
		}
		providers = append(providers, newProvider)
	}
	return providers, nil
}

func VerifyCreds(ctx context.Context, config aws.Config) (*sts.GetCallerIdentityOutput, error) {
	stsClient := sts.NewFromConfig(config)
	input := &sts.GetCallerIdentityInput{}
	creds, err := config.Credentials.Retrieve(ctx)
	if err != nil || !creds.HasKeys() {
		return nil, util.AddStackTrace(fmt.Errorf("no AWS credentials found"))
	}
	result, err := stsClient.GetCallerIdentity(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			return nil, util.AddStackTrace(fmt.Errorf(
				"invalid AWS credentials (try running aws sts get-caller-identity). Error code: %v", apiErr.ErrorCode()))
		} else {
			return nil, util.AddStackTrace(err)
		}
	}
	return result, nil
}

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
	"github.com/run-x/cloudgrep/pkg/provider/types"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
	"github.com/run-x/cloudgrep/pkg/util"
	"go.uber.org/zap"
)

type Provider struct {
	config    aws.Config
	isGlobal  bool
	accountId string
}

func (p Provider) String() string {
	realRegion := p.config.Region
	if p.isGlobal {
		realRegion = "Global"
	}
	return fmt.Sprintf("AWS Provider for account %v, region %v", p.accountId, realRegion)
}

func (p Provider) FetchFunctions() map[string]types.FetchFunc {
	funcMap := make(map[string]types.FetchFunc)
	for resourceType, mapping := range p.getTypeMapping() {
		if p.isGlobal != mapping.IsGlobal {
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

	region := p.config.Region
	if p.isGlobal {
		region = "global"
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
	defaultConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	identity, err := VerifyCreds(ctx, defaultConfig)
	if err != nil {
		return nil, err
	}
	logger.Sugar().Infof("Using the following identity: %v", *identity.Arn)

	regions := cfg.Regions
	if len(regions) == 0 {
		regions = []string{"global", defaultConfig.Region}
	}
	logger.Sugar().Infof("Will look in regions %v", regions)
	var providers []types.Provider
	for _, region := range regions {
		newConfig := defaultConfig.Copy()
		if region != "global" {
			newConfig.Region = region
		}
		newProvider := Provider{isGlobal: region == "global", config: newConfig, accountId: *identity.Account}
		providers = append(providers, newProvider)
	}
	return providers, nil
}

func VerifyCreds(ctx context.Context, config aws.Config) (*sts.GetCallerIdentityOutput, error) {
	stsClient := sts.NewFromConfig(config)
	input := &sts.GetCallerIdentityInput{}
	creds, err := config.Credentials.Retrieve(ctx)
	if err != nil || !creds.HasKeys() {
		return nil, util.AddStackStrace(fmt.Errorf("no AWS credentials found"))
	}
	result, err := stsClient.GetCallerIdentity(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			return nil, util.AddStackStrace(fmt.Errorf(
				"invalid AWS credentials (try running aws sts get-caller-identity). Error code: %v", apiErr.ErrorCode()))
		} else {
			return nil, util.AddStackStrace(err)
		}
	}
	return result, nil
}

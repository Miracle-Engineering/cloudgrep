package regions

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

const Global = "global"
const All = "all"

func SelectRegions(ctx context.Context, configuredRegions []string, awsConfig aws.Config) ([]Region, error) {
	var err error

	if len(configuredRegions) == 1 && configuredRegions[0] == All {
		return allRegions(ctx, awsConfig)
	}

	if len(configuredRegions) > 0 {
		// If regions were configured, use those
		err = validateRegions(configuredRegions)
		if err != nil {
			return nil, fmt.Errorf("unable to use configured regions: %w", err)
		}

		return regionsFromStrings(configuredRegions), nil
	}

	region := awsConfig.Region

	// If we can't detect region automatically, prompt for it
	if region == "" {
		region, err = promptForRegion(ctx)
		if err != nil {
			if err == ctx.Err() {
				return nil, err
			}

			return nil, fmt.Errorf("error prompting for region: %w", err)
		}
	} else {
		err = validateRegions([]string{region})
		if err != nil {
			return nil, err
		}
	}

	regions := []string{Global, region}

	// Always include global region without explicit configuration excluding it
	return regionsFromStrings(regions), err
}

func IsValid(region string) bool {
	if region == Global {
		return true
	}

	_, has := officialRegions[region]
	return has
}

func ConfigureConfigRegion(cfg *aws.Config, regions []Region) {
	if cfg == nil {
		panic("unexpected nil cfg")
	}

	if cfg.Region != "" {
		return
	}

	cfg.Region = "us-east-1"
	for _, region := range regions {
		if !region.IsGlobal() {
			cfg.Region = region.ID()
			return
		}
	}
}

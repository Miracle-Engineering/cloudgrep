package regions

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/manifoldco/promptui"
)

var officialRegions map[string]endpoints.Region

func init() {
	partition := endpoints.AwsPartition()
	officialRegions = partition.Regions()
}

func promptForRegion(ctx context.Context) (string, error) {
	validate := func(input string) error {
		if !IsValid(input) {
			return fmt.Errorf("invalid AWS region code: %v please refer to https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-available-regions", input)
		}
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		prompt := promptui.Prompt{
			Label:    "No default AWS region found, please specify one region code",
			Validate: validate,
		}

		result, err := prompt.Run()

		if err != nil {
			if err == promptui.ErrInterrupt {
				os.Exit(1)
			}

			fmt.Printf("Encountered issue with input: %v\nPlease try again", err)
		} else {
			return result, nil
		}
	}
}

func validateRegions(regions []string) error {
	var badRegions []string
	for _, region := range regions {
		if IsValid(region) {
			continue
		}

		badRegions = append(badRegions, region)
	}

	if len(badRegions) == 0 {
		return nil
	}

	plural := "regions"
	if len(badRegions) == 1 {
		plural = "region"
	}

	return fmt.Errorf("invalid AWS %s: %s", plural, strings.Join(badRegions, ", "))
}

func listAvailableRegions(ctx context.Context, cfg aws.Config) ([]string, error) {
	cfg = cfg.Copy()

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeRegionsInput{}

	output, err := client.DescribeRegions(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to call ec2:DescribeRegions: %w", err)
	}

	var regions []string
	for _, region := range output.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}

func allRegions(ctx context.Context, cfg aws.Config) ([]Region, error) {
	availableRegions, err := listAvailableRegions(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot get all regions: %w", err)
	}

	regions := make([]Region, 0, len(officialRegions)+1)
	regions = append(regions, Region{})

	for _, regionName := range availableRegions {
		region, has := officialRegions[regionName]
		if !has {
			continue
		}

		regions = append(regions, Region{region: &region})
	}

	return regions, nil
}

// func allRegions() []Region {
// 	regions := make([]Region, 0, len(officialRegions)+1)

// 	regions = append(regions, Region{})

// 	for _, officialRegion := range officialRegions {
// 		region := officialRegion
// 		regions = append(regions, Region{region: &region})
// 	}

// 	return regions
// }

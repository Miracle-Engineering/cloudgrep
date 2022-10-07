package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter"
)

// displayIdArn modifies the resource's display ID to be the "resource" part of the ARN.
// Errors if the effective display ID is not an ARN.
func displayIdArn(ctx context.Context, resource *model.Resource) error {
	displayId := resource.EffectiveDisplayId()

	parsed, err := arn.Parse(displayId)
	if err != nil {
		return fmt.Errorf("effective display ID ('%s') is not a valid AWS ARN: %w", displayId, err)
	}

	resource.DisplayId = parsed.Resource
	return nil
}

// displayIdArnPrefix returns a transform func to modify the resource's display ID to be the "resource" part of the ARN, and additionally
// drop some specified prefix from the start of the new display ID.
// Errors if the effective display ID is not an ARN or the specified prefix is not present.
func displayIdArnPrefix(prefix string) resourceconverter.TransformResourceFunc {
	return func(ctx context.Context, resource *model.Resource) error {
		displayId := resource.EffectiveDisplayId()

		parsed, err := arn.Parse(displayId)
		if err != nil {
			return fmt.Errorf("effective display ID ('%s') is not a valid AWS ARN: %w", displayId, err)
		}

		arnResource := parsed.Resource
		if !strings.HasPrefix(arnResource, prefix) {
			return fmt.Errorf("effective display ID ('%s') ARN is missing prefix: %s", displayId, prefix)
		}

		arnResource = strings.TrimPrefix(arnResource, prefix)

		resource.DisplayId = arnResource
		return nil
	}
}

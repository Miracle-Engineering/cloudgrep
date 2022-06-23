package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/run-x/cloudgrep/pkg/util"
)

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

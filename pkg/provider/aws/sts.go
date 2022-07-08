package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

//assumeRoleCredentials will assume the role and returns the credentials to use
func assumeRoleCredentials(config aws.Config, roleArn string) (aws.Credentials, error) {
	//generate a unique session name
	sessionName := fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	response, err := sts.NewFromConfig(config).AssumeRole(context.TODO(), &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
	})
	if err != nil {
		return aws.Credentials{}, err
	}
	return aws.Credentials{
		AccessKeyID:     *response.Credentials.AccessKeyId,
		SecretAccessKey: *response.Credentials.SecretAccessKey,
		SessionToken:    *response.Credentials.SessionToken,
		CanExpire:       true,
		Expires:         *response.Credentials.Expiration,
	}, nil
}

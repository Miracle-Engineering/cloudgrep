package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/mocks/pkg/provider/aws"
)

func TestEC2(t *testing.T) {
	provider := AWSProvider{}
	mock := aws.NewEC2Client(t)
	provider.ec2Client = mock
}

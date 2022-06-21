package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func listHealthCheckTagsInput(input *route53.ListTagsForResourcesInput) error {
	input.ResourceType = types.TagResourceTypeHealthcheck
	for idx, id := range input.ResourceIds {
		id = strings.TrimPrefix(id, "/healthcheck/")
		input.ResourceIds[idx] = id
	}

	return nil
}

func listHostedZoneTagsInput(input *route53.ListTagsForResourcesInput) error {
	input.ResourceType = types.TagResourceTypeHostedzone
	for idx, id := range input.ResourceIds {
		id = strings.TrimPrefix(id, "/hostedzone/")
		input.ResourceIds[idx] = id
	}

	return nil
}

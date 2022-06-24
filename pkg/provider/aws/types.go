package aws

import (
	"github.com/run-x/cloudgrep/pkg/provider/types"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

type mapper struct {
	IdField         string
	TagField        resourceconverter.TagField
	FetchFunc       types.FetchFunc
	IsGlobal        bool
	UseMapConverter bool
	// ServiceEndpointID is the identifier for the service in the endpoints file located at https://github.com/aws/aws-sdk-go/blob/v1.44.33/aws/endpoints/defaults.go.
	// For example, for the `elb` service, the EndpointID is `elasticloadbalancing`.
	// Used to detect regions that don't support specific services.
	// If not set, it is assumed this type's service is supported in all regions.
	ServiceEndpointID string
}

func (p *Provider) getTypeMapping() map[string]mapper {
	return p.buildTypeMapping()
}

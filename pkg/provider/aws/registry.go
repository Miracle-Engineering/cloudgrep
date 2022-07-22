package aws

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func (p *Provider) buildTypeMapping() map[string]mapper {
	mapping := map[string]mapper{}

	p.registerGeneratedTypes(mapping)
	p.register_s3(mapping)
	p.register_eks(mapping)
	p.register_cloudfront(mapping)
	p.register_sqs(mapping)
	p.register_iam_manual(mapping)
	return mapping
}

// SupportedResources returns the resources that are supported by the AWS provider.
func SupportedResources() []string {
	p := Provider{}

	resources := maps.Keys(p.buildTypeMapping())
	slices.Sort(resources)
	return resources
}

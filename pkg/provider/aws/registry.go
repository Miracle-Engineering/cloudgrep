package aws

func (p *Provider) buildTypeMapping() map[string]mapper {
	mapping := map[string]mapper{}

	p.registerGeneratedTypes(mapping)
	p.register_s3(mapping)
	p.register_eks(mapping)

	return mapping
}

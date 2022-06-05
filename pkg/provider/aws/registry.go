package aws

func (p *Provider) buildTypeMapping() map[string]mapper {
	mapping := map[string]mapper{}

	p.registerAllTypes(mapping)

	return mapping
}

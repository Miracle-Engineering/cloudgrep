package aws

import ()

func (p *Provider) registerAllTypes(mapping map[string]mapper) {
	p.register_ec2(mapping)
	p.register_rds(mapping)
}

package aws

import ()

func (p *Provider) registerGeneratedTypes(mapping map[string]mapper) {
	p.registerAutoscaling(mapping)
	p.registerEc2(mapping)
	p.registerElasticache(mapping)
	p.registerElb(mapping)
	p.registerIam(mapping)
	p.registerLambda(mapping)
	p.registerRds(mapping)
	p.registerRoute53(mapping)
	p.registerSns(mapping)
}

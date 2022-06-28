package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider/types"
)

// FakeProvider provides mechanisms to test systems that require providers
// Has facilities for variable resource counts, returning errors, and artificial delay
type FakeProvider struct {
	// ID is used for prefixes on resource types and identifiers
	ID string

	Foo FakeProviderResourceConfig
	Bar FakeProviderResourceConfig
}

var _ types.Provider = &FakeProvider{}
var _ fmt.Stringer = &FakeProvider{}

func (p *FakeProvider) Id() string {
	return p.String()
}

func (p *FakeProvider) String() string {
	if p.ID == "" {
		return "fake"
	}

	return p.ID
}

func (p *FakeProvider) types() map[string]*FakeProviderResourceConfig {
	return map[string]*FakeProviderResourceConfig{
		"Bar": &p.Bar,
		"Foo": &p.Foo,
	}
}

func (p *FakeProvider) FetchFunctions() map[string]types.FetchFunc {
	mapping := make(map[string]types.FetchFunc)
	for typ, config := range p.types() {
		resourceType := fmt.Sprintf("%s.%s", p.String(), typ)

		mapping[resourceType] = p.makeFetchFunc(resourceType, config)
	}

	return mapping
}

func (p *FakeProvider) makeFetchFunc(typ string, config *FakeProviderResourceConfig) types.FetchFunc {
	return func(ctx context.Context, output chan<- model.Resource) error {
		defer config.addRun()

		resources := make([]model.Resource, config.Count)

		p.annotateResources(typ, resources)

		for _, resource := range resources {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(config.DelayBefore):
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- resource:
			}

		}

		return config.ErrorAfterCount
	}
}

func (p *FakeProvider) annotateResources(typ string, resources []model.Resource) {
	for idx, resource := range resources {
		id := fmt.Sprintf("%s-%s-%d", p.String(), typ, idx)
		resource.Id = id
		resources[idx] = resource
	}
}

func (p *FakeProvider) TotalRuns() int {
	var count int
	for _, config := range p.types() {
		count += config.RunCount()
	}

	return count
}

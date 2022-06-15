package aws

import (
	"github.com/run-x/cloudgrep/pkg/provider/types"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

type mapper struct {
	IdField   string
	TagField  resourceconverter.TagField
	FetchFunc types.FetchFunc
	IsGlobal  bool
}

func (p *Provider) getTypeMapping() map[string]mapper {
	return p.buildTypeMapping()
}

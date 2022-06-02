package sequencer

import (
	"context"

	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/provider"
)

type Sequencer interface {
	Run(ctx context.Context, ds datastore.Datastore, providers []provider.Provider) error
}

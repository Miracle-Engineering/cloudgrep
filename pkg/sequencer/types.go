package sequencer

import (
	"context"

	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/provider2"
)

type Sequencer interface {
	Run(ctx context.Context, ds datastore.Datastore, providers []provider2.Provider) error
}

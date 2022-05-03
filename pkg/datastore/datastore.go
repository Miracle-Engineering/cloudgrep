package datastore

import (
	"context"
	"fmt"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
)

//Datastore provides an interface to read/write/update to a store
type Datastore interface {
	init(context.Context, config.Config) error
	GetResources(context.Context, model.Filter) ([]*model.Resource, error)
	WriteResources(context.Context, []*model.Resource) error
}

func NewDatastore(ctx context.Context, cfg config.Config) (Datastore, error) {
	var datastore Datastore
	if cfg.Datastore.Type == "memory" {
		datastore = &MemoryStore{}
	}
	if datastore == nil {
		return nil, fmt.Errorf("unknown datastore type '%v'", cfg.Datastore.Type)
	}
	cfg.Logging.Logger.Sugar().Infow("Creating a datastore",
		zap.String("type", cfg.Datastore.Type))
	err := datastore.init(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return datastore, nil
}

package datastore

import (
	"context"
	"fmt"

	"github.com/juandiegopalomino/cloudgrep/pkg/config"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"go.uber.org/zap"
)

//Datastore provides an interface to read/write/update to a store
type Datastore interface {
	GetResource(context.Context, string) (*model.Resource, error)
	GetResources(context.Context, []byte) (model.ResourcesResponse, error)
	WriteResources(context.Context, model.Resources) error
	Stats(context.Context) (model.Stats, error)
	WriteEvent(context.Context, model.Event) error
	EngineStatus(ctx context.Context) (model.Event, error)
	Ping() error
}

func NewDatastore(ctx context.Context, cfg config.Config, logger *zap.Logger) (Datastore, error) {
	logger.Sugar().Infow("Creating a datastore", zap.String("type", cfg.Datastore.Type))
	switch cfg.Datastore.Type {
	case "sqlite":
		return NewSQLiteStore(ctx, cfg, logger)
	}
	return nil, fmt.Errorf("unknown datastore type '%v'", cfg.Datastore.Type)
}

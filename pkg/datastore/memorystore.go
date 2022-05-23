package datastore

import (
	"context"
	"errors"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
)

//MemoryStore stores the last resources written in memory and return it without any filtering
//Not for production use!
type MemoryStore struct {
	logger    *zap.Logger
	resources []*model.Resource
}

func NewMemoryStore(ctx context.Context, cfg config.Config) *MemoryStore {
	if !cfg.Logging.IsDev() {
		cfg.Logging.Logger.Warn("MemoryStore should not be used for production")
	}
	datastore := MemoryStore{}
	datastore.logger = cfg.Logging.Logger
	return &datastore
}

func (m *MemoryStore) GetResource(ctx context.Context, id string) (*model.Resource, error) {
	for _, r := range m.resources {
		if r.Id == id {
			return r, nil
		}
	}
	//not found
	return nil, nil
}

func (m *MemoryStore) GetResources(ctx context.Context, filter model.Filter) ([]*model.Resource, error) {
	result := m.resources
	if !filter.IsEmpty() {
		return nil, errors.New("not implemented")
	}
	m.logger.Sugar().Infow("Getting resources: ",
		zap.Object("filter", filter),
		zap.Int("count", len(result)),
	)
	return result, nil
}

func (m *MemoryStore) WriteResources(ctx context.Context, resources []*model.Resource) error {
	m.logger.Sugar().Infow("Writting resources: ",
		zap.Int("count", len(resources)),
	)
	m.resources = resources
	return nil
}

func (m *MemoryStore) Stats(context.Context) (model.Stats, error) {
	return model.Stats{}, errors.New("not implemented")
}

func (m *MemoryStore) GetFields(context.Context) (model.Fields, error) {
	return nil, errors.New("not implemented")
}

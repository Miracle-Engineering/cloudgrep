package datastore

import (
	"context"

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

func (m *MemoryStore) init(ctx context.Context, cfg config.Config) error {
	if !cfg.Logging.IsDev() {
		cfg.Logging.Logger.Warn("MemoryStore should not be used for production")
	}
	m.logger = cfg.Logging.Logger
	return nil
}

func (m *MemoryStore) GetResources(ctx context.Context, filter model.Filter) ([]*model.Resource, error) {
	result := m.resources
	m.logger.Sugar().Infow("Getting resources: ",
		zap.String("filter", filter.String()),
		zap.Int("count", len(result)),
	)
	return result, nil
}

func (m *MemoryStore) WriteResources(ctx context.Context, resources []*model.Resource) error {
	m.resources = resources
	m.logger.Sugar().Infow("Writting resources: ",
		zap.Int("count", len(resources)),
	)
	return nil
}

package datastore

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLiteStore struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewSQLiteStore(ctx context.Context, cfg config.Config) (*SQLiteStore, error) {
	s := SQLiteStore{}
	logLevel := logger.Error
	if cfg.Logging.IsDev() {
		//log all SQL queries
		logLevel = logger.Info
	}
	//gormLogger has it's own logger for SQL queries - better than zaplog for that purpose
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	s.logger = cfg.Logging.Logger
	var err error
	//create the DB client
	if s.db, err = gorm.Open(sqlite.Open(cfg.Datastore.DataSourceName),
		&gorm.Config{Logger: gormLogger}); err != nil {
		return nil, fmt.Errorf("can't create the SQLite database: %w", err)
	}

	// Migrate the schema
	if err = s.db.AutoMigrate(&model.Resource{}, &model.Property{}, &model.Tag{}); err != nil {
		return nil, fmt.Errorf("can't create the SQLite data model: %w", err)
	}

	return &s, nil
}

func (s *SQLiteStore) GetResources(ctx context.Context, filter model.Filter) ([]*model.Resource, error) {
	var resources []*model.Resource
	result := s.db.Preload("Tags").Preload("Properties").Find(&resources)
	s.logger.Sugar().Infow("Getting resources: ",
		zap.String("filter", filter.String()),
		zap.Int("count", int(result.RowsAffected)),
	)
	if result.Error != nil {
		return nil, fmt.Errorf("can't get resources from database: %w", result.Error)
	}
	return resources, nil
}

func (s *SQLiteStore) WriteResources(ctx context.Context, resources []*model.Resource) error {
	result := s.db.Create(resources)
	s.logger.Sugar().Infow("Writting resources: ",
		zap.Int("count", int(result.RowsAffected)),
	)
	if result.Error != nil {
		return fmt.Errorf("can't write resources to database: %w", result.Error)
	}
	return nil
}

package datastore

import (
	"context"
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
	s.db, err = gorm.Open(sqlite.Open(cfg.Datastore.DataSourceName), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	err = s.db.AutoMigrate(&model.Resource{}, &model.Property{}, &model.Tag{})
	if err != nil {
		return nil, err
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
	return resources, nil
}

func (s *SQLiteStore) WriteResources(ctx context.Context, resources []*model.Resource) error {
	result := s.db.Create(resources)
	s.logger.Sugar().Infow("Writting resources: ",
		zap.Int("count", int(result.RowsAffected)),
	)
	return result.Error
}

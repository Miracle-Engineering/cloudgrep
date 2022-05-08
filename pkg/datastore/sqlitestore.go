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

type resourceId string

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
	//create the DB client
	var err error
	s.db, err = gorm.Open(sqlite.Open(cfg.Datastore.DataSourceName),
		&gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("can't create the SQLite database: %w", err)
	}

	// Migrate the schema
	if err = s.db.AutoMigrate(&model.Resource{}, &model.Property{}, &model.Tag{}); err != nil {
		return nil, fmt.Errorf("can't create the SQLite data model: %w", err)
	}

	return &s, nil
}

func (s *SQLiteStore) getResourceIds(ctx context.Context, filter model.Filter) ([]resourceId, error) {
	var resourceIds []resourceId
	db := s.db.Model(&model.Tag{}).Select("resource_id")
	if !filter.IsEmpty() {
		for _, tag := range filter.Tags {
			if tag.Value == "*" {
				//wild card - check for any value
				db = db.Or("key=? and value like ?", tag.Key, "%")
			} else {
				//exact value
				db = db.Or("key=? and value=?", tag.Key, tag.Value)
			}
		}
		// our data model has one row per tag, we need a group by to only keep the results matching all tags
		/* table tags
		//original data, tag filter = team=infra and cluster=staging
		1. i-123	cluster		staging
		2. i-123	team		infra
		3. i-124	cluster		staging
		4. i-124	team		dev

		//where team=infra OR cluster=staging
		1. i-123	cluster		staging
		2. i-123	team		infra
		3. i-124	cluster		staging

		// group by + count on id
		1. i-123	2 <- we only want to keep this row
		2. i-124	1
		*/
		db = db.Group("resource_id").
			Having("count(resource_id)=?", filter.Tags.DistinctKeys())
	}

	result := db.Find(&resourceIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return resourceIds, nil
}

func (s *SQLiteStore) getResourcesById(ctx context.Context, ids []resourceId) ([]*model.Resource, error) {
	var resources []*model.Resource
	if len(ids) == 0 {
		return resources, nil
	}
	result := s.db.Preload("Tags").Preload("Properties").Find(&resources, ids)

	if result.Error != nil {
		return nil, result.Error
	}
	return resources, nil

}

func (s *SQLiteStore) GetResources(ctx context.Context, filter model.Filter) ([]*model.Resource, error) {
	var resources []*model.Resource

	ids, err := s.getResourceIds(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("can't get resources from database: %w", err)
	}
	resources, err = s.getResourcesById(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("can't get resources from database: %w", err)
	}
	s.logger.Sugar().Infow("Getting resources: ",
		zap.Object("filter", filter),
		zap.Int("count", len(resources)),
	)
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

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
	db     *CloudGrepDB
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
	//create the DB client
	db, err := gorm.Open(sqlite.Open(cfg.Datastore.DataSourceName),
		&gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("can't create the SQLite database: %w", err)
	}
	s.db = &CloudGrepDB{db}

	// Migrate the schema
	if err = s.db.AutoMigrate(&model.Resource{}, &model.Property{}, &model.Tag{}); err != nil {
		return nil, fmt.Errorf("can't create the SQLite data model: %w", err)
	}

	return &s, nil
}

type CloudGrepDB struct {
	*gorm.DB
}

//joinOnTags joins on tags table using resource_id as key
func (db *CloudGrepDB) joinOnTags() *CloudGrepDB {
	return &CloudGrepDB{db.Joins("JOIN tags on tags.resource_id=resources.id")}
}

//filter build the query to only return filtered results
func (db *CloudGrepDB) filter(filter model.Filter) *CloudGrepDB {
	if filter.IsEmpty() {
		return db
	}
	//working a session allow chaining using OR on where statement
	// without session: (=tag1) OR ((=tag1) AND (=tag2))
	// with session: (=tag1) OR (=tag2)
	table := db.joinOnTags().Session(&gorm.Session{})

	// use a group condition to keep all the matching tags
	// where (=tag1) OR (=tag2)
	// https://gorm.io/docs/advanced_query.html#Group-Conditions
	result := table
	for _, tag := range filter.Tags {
		if tag.Value == "*" {
			//wild card - check for any value
			result = result.
				Or(
					table.Where("tags.key=? and tags.value like ?", tag.Key, "%"),
				)
		} else {
			//exact value
			result = result.
				Or(
					table.Where("tags.key=? and tags.value=?", tag.Key, tag.Value),
				)
		}
	}
	// our data model has one row per tag, we need a group by to only keep the results matching all tags
	/* table resource join tag
	1. i-123	us-east-1	cluster		staging
	2. i-123	us-east-1	team		infra
	3. i-124	us-east-1	cluster		staging
	3. i-124	us-east-1	team		dev

	when searching on team=infra and cluster=staging, we should only return i-123
	# group by + count on id
	1. i-123	2 <- we only want to keep this row
	2. i-124	1

	*/
	result = result.Group("resources.id").
		Having("count(resources.id)=?", filter.Tags.DistinctKeys())

	return &CloudGrepDB{result}
}

func (s *SQLiteStore) GetResources(ctx context.Context, filter model.Filter) ([]*model.Resource, error) {
	var resources []*model.Resource

	result := s.db.filter(filter).
		Preload("Tags").Preload("Properties").
		Find(&resources)

	if result.Error != nil {
		return nil, fmt.Errorf("can't get resources from database: %w", result.Error)
	}
	s.logger.Sugar().Infow("Getting resources: ",
		zap.Object("filter", filter),
		zap.Int("count", int(result.RowsAffected)),
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

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

func (s *SQLiteStore) getAllResourceIds(ctx context.Context) ([]model.ResourceId, error) {
	var resourceIds []model.ResourceId
	result := s.db.Model(&model.Resource{}).Select("id").Distinct().Find(&resourceIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return resourceIds, nil
}

func (s *SQLiteStore) getResourceIdsByTag(ctx context.Context, tags model.Tags) ([]model.ResourceId, error) {
	var resourceIds []model.ResourceId
	db := s.db.Model(&model.Tag{}).Select("resource_id").Distinct()
	if !tags.Empty() {
		for _, tag := range tags {
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
			Having("count(resource_id)=?", tags.DistinctKeys())
	}
	result := db.Find(&resourceIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return resourceIds, nil
}

//getResourceIdsForFilter apply the resource filter and return the ids to include, ids to exclude and potential error
func (s *SQLiteStore) getResourceIdsForFilter(ctx context.Context, filter model.Filter) ([]model.ResourceId, []model.ResourceId, error) {

	//apply filter tags:include
	includeTags := filter.TagsInclude()
	var err error
	var idsInclude []model.ResourceId
	if includeTags.Empty() {
		//no include filter set - include all
		idsInclude, err = s.getAllResourceIds(ctx)
	} else {
		idsInclude, err = s.getResourceIdsByTag(ctx, includeTags)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("can't get resources from database: %w", err)
	}

	//apply filter tags:exclude
	var idsExclude []model.ResourceId
	if tagsExclude := filter.TagsExclude(); !tagsExclude.Empty() {
		idsExclude, err = s.getResourceIdsByTag(ctx, filter.TagsExclude())
		if err != nil {
			return nil, nil, fmt.Errorf("can't get resources from database: %w", err)
		}
	}
	return idsInclude, idsExclude, nil
}

func (s *SQLiteStore) getResourcesById(ctx context.Context, idsInclude []model.ResourceId, idsExclude []model.ResourceId) ([]*model.Resource, error) {
	var resources []*model.Resource
	if len(idsInclude) == 0 {
		return resources, nil
	}
	db := whereResourceIds(s.db, "id", idsInclude, idsExclude)
	db = db.Preload("Tags").Preload("Properties").Find(&resources)

	if db.Error != nil {
		return nil, db.Error
	}
	return resources, nil

}
func (s *SQLiteStore) GetResource(ctx context.Context, id model.ResourceId) (*model.Resource, error) {
	resources, err := s.getResourcesById(ctx, []model.ResourceId{model.ResourceId(id)}, nil)
	if err != nil {
		return nil, fmt.Errorf("can't get resource from database: %w", err)
	}
	if len(resources) == 0 {
		//not found
		return nil, nil
	}
	s.logger.Sugar().Infow("Getting resource: ",
		zap.String("id", string(id)),
	)
	return resources[0], nil
}
func (s *SQLiteStore) GetResources(ctx context.Context, filter model.Filter) ([]*model.Resource, error) {
	//get resource ids for current filter
	idsInclude, idsExclude, err := s.getResourceIdsForFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	var resources []*model.Resource
	resources, err = s.getResourcesById(ctx, idsInclude, idsExclude)
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

func (s *SQLiteStore) Stats(context.Context) (model.Stats, error) {
	var count int64
	result := s.db.Table("resources").Count(&count)
	if result.Error != nil {
		return model.Stats{}, fmt.Errorf("can't read resources count from database: %w", result.Error)
	}
	return model.Stats{ResourcesCount: int(count)}, nil
}

//getTagValues return the tag values for a given key
func (s *SQLiteStore) getTagValues(ctx context.Context, key string, resourceIdsInclude []model.ResourceId, resourceIdsExclude []model.ResourceId) ([]string, error) {
	/*
		select distinct(value) from tags
		where tags.resource_id in (
			resourceIdsInclude
		)
		and tags.resource_id not in (
			resourceIdsExclude
		)
		and key='?'
	*/
	db := s.db.Model(&model.Tag{}).Select("value").Distinct().
		Where("key=?", key)
	db = whereResourceIds(db, "resource_id", resourceIdsInclude, resourceIdsExclude)

	var values []string
	db = db.Find(&values)
	if db.Error != nil {
		return nil, db.Error
	}
	return values, nil
}

//getTagValues return the tag resource ids for a given key
func (s *SQLiteStore) getTagResourceIds(ctx context.Context, key string, resourceIdsInclude []model.ResourceId, resourceIdsExclude []model.ResourceId) ([]model.ResourceId, error) {
	/*
		select distinct(resource_id) from tags
		where tags.resource_id in (
			resourceIdsInclude
		)
		and tags.resource_id not in (
			resourceIdsExclude
		)
		and key='?'
	*/
	db := s.db.Model(&model.Tag{}).Select("resource_id").Distinct().
		Where("key=?", key)
	db = whereResourceIds(db, "resource_id", resourceIdsInclude, resourceIdsExclude)

	var resource_ids []model.ResourceId
	if len(resourceIdsInclude) == 0 {
		return resource_ids, nil
	} else {
		db = db.Where("resource_id in ?", resourceIdsInclude)
	}
	if len(resourceIdsExclude) != 0 {
		// exlucde these ids
		db = db.Where("resource_id not in ?", resourceIdsExclude)
	}
	db = db.Find(&resource_ids)
	if db.Error != nil {
		return nil, db.Error
	}
	return resource_ids, nil
}

func (s *SQLiteStore) GetTags(ctx context.Context, filter model.Filter, limit int) ([]*model.TagInfo, error) {
	//get resource ids for current filter
	resourceIdsInclude, resourceIdsExclude, err := s.getResourceIdsForFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(resourceIdsInclude) == 0 {
		//current filter has no result
		return nil, nil
	}

	db := s.db.Model(&model.Tag{}).Select("key", "count() as count")
	db = whereResourceIds(db, "resource_id", resourceIdsInclude, resourceIdsExclude)

	var tagInfos []*model.TagInfo
	//get the tags with the highest count, apply the limit
	db = db.
		Group("key").
		Order("count DESC").
		Limit(limit).
		Find(&tagInfos)

	if db.Error != nil {
		return nil, fmt.Errorf("can't read tags from database: %w", db.Error)
	}

	//get the values for each key
	for _, tagInfo := range tagInfos {
		values, err := s.getTagValues(ctx, tagInfo.Key, resourceIdsInclude, resourceIdsExclude)
		if err != nil {
			return nil, err
		}
		resourceIds, err := s.getTagResourceIds(ctx, tagInfo.Key, resourceIdsInclude, resourceIdsExclude)
		if err != nil {
			return nil, err
		}
		tagInfo.Values = values
		tagInfo.ResourceIds = resourceIds
	}

	return tagInfos, nil
}

//whereIds add the where clause to select the provided resource ids
func whereResourceIds(db *gorm.DB, idName string, idsInclude []model.ResourceId, idsExclude []model.ResourceId) *gorm.DB {
	db = db.Where(idName+" in ?", idsInclude)
	if len(idsExclude) != 0 {
		// exlucde these ids
		db = db.Where(idName+" not in ?", idsExclude)
	}
	return db
}

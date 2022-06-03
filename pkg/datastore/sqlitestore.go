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
	logger  *zap.Logger
	db      *gorm.DB
	indexer resourceIndexer
}

type resourceId string

func NewSQLiteStore(ctx context.Context, cfg config.Config, zapLogger *zap.Logger) (*SQLiteStore, error) {
	s := SQLiteStore{}
	logLevel := logger.Error
	if zapLogger.Core().Enabled(zap.DebugLevel) {
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
	s.logger = zapLogger
	//create the DB client
	var err error
	s.db, err = gorm.Open(sqlite.Open(cfg.Datastore.DataSourceName),
		&gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("can't create the SQLite database: %w", err)
	}

	// Migrate the schema
	if err = s.db.AutoMigrate(&model.Resource{}, &model.Tag{}, &model.EngineStatus{}); err != nil {
		return nil, fmt.Errorf("can't create the SQLite data model: %w", err)
	}

	//create the indexer
	s.indexer, err = newResourceIndexer(ctx, s.logger, s.db)
	if err != nil {
		return nil, fmt.Errorf("can't create the query builder: %w", err)
	}

	return &s, nil
}

func (s *SQLiteStore) Ping() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (s *SQLiteStore) getResourcesById(ctx context.Context, ids []resourceId) ([]*model.Resource, error) {
	var resources []*model.Resource
	if len(ids) == 0 {
		return resources, nil
	}
	db := s.db.Preload("Tags").Find(&resources, ids)

	if db.Error != nil {
		return nil, db.Error
	}
	return resources, nil

}
func (s *SQLiteStore) GetResource(ctx context.Context, id string) (*model.Resource, error) {
	resources, err := s.getResourcesById(ctx, []resourceId{resourceId(id)})
	if err != nil {
		return nil, fmt.Errorf("can't get resource from database: %w", err)
	}
	if len(resources) == 0 {
		//not found
		return nil, nil
	}
	s.logger.Sugar().Infow("Getting resource: ",
		zap.String("id", id),
	)
	return resources[0], nil
}

func (s *SQLiteStore) WriteResources(ctx context.Context, resources []*model.Resource) error {
	if len(resources) == 0 {
		//nothing to write
		return nil
	}
	result := s.db.Create(resources)
	s.logger.Sugar().Infow("Writting resources: ",
		zap.Int("count", int(result.RowsAffected)),
	)

	if result.Error != nil {
		return fmt.Errorf("can't write resources to database: %w", result.Error)
	}
	if err := s.indexer.writeResourceIndexes(ctx, s.db, resources); err != nil {
		return err
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

func (s *SQLiteStore) getResourceField(name string) (model.Field, error) {
	/*
		SELECT DISTINCT `type` , count(*) as count
		FROM `resources`
		group by  `type`
		order by `type`
		sort by `count` desc
	*/
	rows, err := s.db.Model(&model.Resource{}).Select(name, "count() as count").
		Distinct().
		Group(name).
		Order("count desc").
		Rows()
	if err != nil {
		return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
	}
	defer rows.Close()
	field := model.Field{
		Name: name,
	}
	var totalCount int
	for rows.Next() {
		var value string
		var count int
		err = rows.Scan(&value, &count)
		if err != nil {
			return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
		}
		field.Values = append(field.Values, model.FieldValue{
			Value: value,
			Count: count,
		})
		totalCount = totalCount + count
	}
	field.Count = totalCount
	return field, nil
}

func (s *SQLiteStore) getTagFields() (model.Fields, error) {
	fields, err := s.getTagKeys()
	if err != nil {
		return model.Fields{}, fmt.Errorf("can't get tag keys from database: %w", err)
	}
	var result model.Fields
	for _, f := range fields {
		values, err := s.getTagValues(f.Name)
		if err != nil {
			return nil, err
		}
		f.Values = values
		result = append(result, f)
	}
	return result, nil
}

func (s *SQLiteStore) getTagKeys() (model.Fields, error) {
	/*
		SELECT distinct(`key`), count() as count
		FROM `tags`
		group by key
		order by count desc
	*/
	rows, err := s.db.Model(&model.Tag{}).Select("key", "count() as count").
		Distinct().
		Group("key").
		Order("count desc").
		Rows()
	if err != nil {
		return model.Fields{}, fmt.Errorf("can't get tag keys from database: %w", err)
	}
	var fields model.Fields
	defer rows.Close()
	for rows.Next() {
		var key string
		var count int
		err = rows.Scan(&key, &count)
		if err != nil {
			return nil, err
		}
		field := model.Field{
			Name:  key,
			Count: count,
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func (s *SQLiteStore) getTagValues(key string) (model.FieldValues, error) {
	/*
		SELECT distinct(`value`), count() as count
		FROM `tags`
		where key=?
		group by value
		order by count desc
	*/
	var values model.FieldValues
	db := s.db.Model(&model.Tag{}).Select("value", "count() as count").
		Distinct().
		Where("key=?", key).
		Group("value").
		Order("count desc").
		Find(&values)

	if db.Error != nil {
		return nil, fmt.Errorf("can't get tag value for key '%v' : %w", key, db.Error)
	}
	return values, nil
}

func (s *SQLiteStore) GetFields(context.Context) (model.FieldGroups, error) {
	var fieldGroups model.FieldGroups

	//get core fields
	coreGroup := model.FieldGroup{
		Name: "core",
	}
	for _, name := range []string{"region", "type"} {
		field, err := s.getResourceField(name)
		if err != nil {
			return nil, err
		}
		coreGroup.Fields = append(coreGroup.Fields, field)
	}
	fieldGroups = append(fieldGroups, coreGroup)

	//get tag fields
	tagFields, err := s.getTagFields()
	if err != nil {
		return nil, err
	}
	tagsGroup := model.FieldGroup{
		Name:   "tags",
		Fields: tagFields,
	}
	fieldGroups = append(fieldGroups, tagsGroup)

	return fieldGroups.AddNullValues(), nil
}

func (s *SQLiteStore) GetResources(ctx context.Context, jsonQuery []byte) ([]*model.Resource, error) {
	ids, err := s.indexer.findResourceIds(*s.db, s.logger, jsonQuery)
	if err != nil {
		return nil, err
	}
	return s.getResourcesById(ctx, ids)
}

func (s *SQLiteStore) WriteEngineStatusStart(ctx context.Context, resource string) error {
	status := model.NewEngineStatus(model.EngineStatusFetching, resource, nil)
	s.logger.Sugar().Infow("Writing Engine Status: ",
		zap.Any("status", status),
	)
	result := s.db.Model(&model.EngineStatus{}).Find(&model.EngineStatus{}, [1]string{resource})
	if result.RowsAffected == 1 {
		result = s.db.Model(&status).Updates(status)
	} else {
		result = s.db.Create(&status)
	}
	if result.Error != nil {
		return fmt.Errorf("can't write engine status to database: %w", result.Error)
	}
	return nil
}

func (s *SQLiteStore) WriteEngineStatusEnd(ctx context.Context, resource string, err error) error {
	var status model.EngineStatus
	if err != nil {
		status = model.NewEngineStatus(model.EngineStatusFailed, resource, err)
	} else {
		status = model.NewEngineStatus(model.EngineStatusSuccess, resource, nil)
	}
	s.logger.Sugar().Infow("Writing Engine Status: ",
		zap.Any("status", status),
	)
	result := s.db.Model(&model.EngineStatus{}).Find(&model.EngineStatus{}, [1]string{resource})
	if result.RowsAffected == 1 {
		result = s.db.Model(&status).Updates(status)
	} else {
		result = s.db.Create(&status)
	}
	if result.Error != nil {
		return fmt.Errorf("can't write engine status to database: %w", result.Error)
	}
	return nil
}

func (s *SQLiteStore) GetEngineStatus(context.Context) (model.EngineStatus, error) {
	var status model.EngineStatus
	db := s.db.Model(&model.EngineStatus{}).Select("*").Order("fetched_at desc").Last(&status)
	if db.Error != nil {
		return model.EngineStatus{}, fmt.Errorf("can't fetch engine status' : %w", db.Error)
	}
	return status, nil
}

package datastore

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type SQLiteStore struct {
	logger  *zap.Logger
	db      *gorm.DB
	indexer resourceIndexer
	//fetchedAt is the last time the resources were fetched
	fetchedAt time.Time
	lock      sync.Mutex
	runId     string
}

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
	s.db, err = gorm.Open(sqlite.Open(s.formatDSN(cfg.Datastore.DataSourceName)),
		&gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("can't create the SQLite database: %w", err)
	}

	// Migrate the schema
	if err = s.db.AutoMigrate(&model.Resource{}, &model.Tag{}, &model.Event{}); err != nil {
		return nil, fmt.Errorf("can't create the SQLite data model: %w", err)
	}

	//create the indexer
	s.indexer, err = newResourceIndexer(ctx, s.logger, s.db)
	if err != nil {
		return nil, fmt.Errorf("can't create the query builder: %w", err)
	}

	return &s, nil
}

func (s *SQLiteStore) formatDSN(dsn string) string {
	if strings.HasPrefix(dsn, "~/") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Errorf("could not handle searching for home directory: %w", err))
		}
		dsn = filepath.Join(dirname, dsn[2:])
	}
	return os.ExpandEnv(dsn)
}

func (s *SQLiteStore) Ping() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (s *SQLiteStore) getResourcesById(ctx context.Context, ids []model.ResourceId) ([]*model.Resource, error) {
	resources := make([]*model.Resource, 0)
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
	resources, err := s.getResourcesById(ctx, []model.ResourceId{model.ResourceId(id)})
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

func (s *SQLiteStore) WriteResources(ctx context.Context, resources model.Resources) error {
	if len(resources) == 0 {
		//nothing to write
		return nil
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	var count int64
	err := s.db.Transaction(func(tx *gorm.DB) error {
		//delete all the previously stored tags if any
		ids := resources.Ids()
		if err := deleteTags(tx, ids); err != nil {
			return err
		}

		// Create or Update all columns
		result := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(resources)
		count = result.RowsAffected

		// Create or Update the resource indexes
		if err := s.indexer.writeResourceIndexes(ctx, tx, resources); err != nil {
			return err
		}
		return result.Error
	})

	s.logger.Sugar().Infow("Writting resources: ", zap.Int64("count", count))

	if err != nil {
		return fmt.Errorf("can't write resources to database: %w", err)
	}
	return nil
}

func resourceCount(db *gorm.DB) (int, error) {
	var count int64
	return int(count), db.Table("resources").Count(&count).Error
}

func (s *SQLiteStore) Stats(context.Context) (model.Stats, error) {
	count, err := resourceCount(s.db)
	if err != nil {
		return model.Stats{}, fmt.Errorf("can't read stats: %w", err)
	}
	return model.Stats{ResourcesCount: count}, nil
}

func (s *SQLiteStore) getResourceField(name string, ids []model.ResourceId) (model.Field, error) {
	/*
		SELECT DISTINCT `type` , count(*) as count
		FROM `resources`
		group by  `type`
		order by `type`
		sort by `count` desc
	*/
	query := s.db.Model(&model.Resource{}).
		Select(name, "count() as count").
		Distinct().
		Group(name).
		Order("count desc")
	rows, err := query.Rows()
	if err != nil {
		return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
	}
	defer rows.Close()
	field := model.Field{
		Name: name,
	}

	//used to update count later
	var totalCount int
	//use a map for efficient access the field value by their value
	fieldValuesMap := make(map[string]*model.FieldValue)

	for rows.Next() {
		var value string
		var count int
		err = rows.Scan(&value, &count)
		if err != nil {
			return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
		}
		fieldValue := model.FieldValue{
			Value: value,
		}
		if len(ids) == 0 {
			//we don't need to go back to the DB
			fieldValue.Count = fmt.Sprint(count)
			totalCount = totalCount + count
		} else {
			//the count would be updated based on filter
			fieldValue.Count = model.CountValueIgnored
		}
		field.Values = append(field.Values, &fieldValue)
		fieldValuesMap[value] = &fieldValue
	}

	//if there is some ids - do another call to update the counts for the current query
	if len(ids) > 0 {
		rows, err = query.Where("id in ?", ids).Rows()
		if err != nil {
			return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
		}
		defer rows.Close()
		for rows.Next() {
			var value string
			var count int
			err = rows.Scan(&value, &count)
			if err != nil {
				return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
			}
			fieldValuesMap[value].Count = fmt.Sprint(count)
			totalCount = totalCount + count
		}
	}

	field.Count = totalCount
	return field, nil
}

//return the list of tag fields sorted by most popular
func (s *SQLiteStore) getTagFields(ids []model.ResourceId) (model.Fields, error) {

	// the tricky part of this function is to always return the same fields containing all the values but with different count
	// the fields are always visible to the user and are ordered by popularity - if we would change this list for every call the UI would look shaky

	// so this is implemented in 3 steps:
	//1. get all tags keys sorted by most frequent first
	//2. get all tags values sorted by most frequent first
	//3. update count to reflect the current query, which relates to subset of resources (ids params)

	// Steps 1 & 2 ensure that all fields and values are returned in the same order
	// Step 3 updates the count values

	//1.  get all tags keys sorted by most frequent first
	rows, err := s.db.Model(&model.Tag{}).Select("key", "count() as count").
		Group("key").
		Order("count desc").
		Rows()
	if err != nil {
		return model.Fields{}, fmt.Errorf("can't get tag keys from database: %w", err)
	}
	var fields model.Fields
	//use a map for efficient access (the model is a slice)
	mapFields := make(map[string]*model.Field)
	defer rows.Close()
	for rows.Next() {
		var key string
		var count int
		err = rows.Scan(&key, &count)
		if err != nil {
			return nil, err
		}
		field := model.Field{
			Name: key,
		}
		if len(ids) == 0 {
			//the current count is for all resources, use it
			field.Count = count
		}
		mapFields[key] = &field
		fields = append(fields, &field)
	}

	//2.  get all tags values sorted by most frequent first
	rows, err = s.db.Raw("SELECT key, value, count() as count FROM tags group by key, value").Rows()
	if err != nil {
		return model.Fields{}, fmt.Errorf("can't get tag values from database: %w", err)
	}
	keyFunc := func(key, val string) string {
		return fmt.Sprintf("%v:%v", key, val)
	}
	//use a map for efficient access
	mapFieldsValue := make(map[string]*model.FieldValue)
	defer rows.Close()
	for rows.Next() {
		var key string
		var value string
		var count int
		err = rows.Scan(&key, &value, &count)
		if err != nil {
			return nil, err
		}
		//add the tag value
		fieldVal := model.FieldValue{
			Value: value,
		}
		if len(ids) == 0 {
			//the current count is for all resources, use it
			fieldVal.Count = fmt.Sprint(count)
		} else {
			//this value will be set using the ids
			fieldVal.Count = "-"
		}
		field := mapFields[key]
		field.Values = append(field.Values, &fieldVal)

		//update map for fast access
		mapFieldsValue[keyFunc(key, value)] = &fieldVal
	}

	if len(ids) == 0 {
		//we are done, no need to get specific values
		return fields, nil
	}

	//3. update count to reflect the current query
	rows, err = s.db.Raw(`SELECT key, value, count() as count `+
		`FROM tags `+
		`Where resource_id in ? `+
		`group by key, value`, ids).Rows()
	if err != nil {
		return model.Fields{}, fmt.Errorf("can't get tag values from database: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var value string
		var count int
		err = rows.Scan(&key, &value, &count)
		if err != nil {
			return nil, err
		}

		//update the tag value
		fieldValue := mapFieldsValue[keyFunc(key, value)]
		fieldValue.Count = fmt.Sprint(count)

		//updte the count of the field
		field := mapFields[key]
		field.Count = field.Count + count
	}

	return fields, nil
}

//TODO remove this api
func (s *SQLiteStore) GetFields(ctx context.Context) (model.FieldGroups, error) {
	return s.getFields(ctx, nil)
}

func (s *SQLiteStore) getFields(ctx context.Context, ids []model.ResourceId) (model.FieldGroups, error) {
	var fieldGroups model.FieldGroups

	//get core fields
	coreGroup := model.FieldGroup{
		Name: model.FieldGroupCore,
	}
	for _, name := range []string{"region", "type"} {
		field, err := s.getResourceField(name, ids)
		if err != nil {
			return nil, err
		}
		coreGroup.Fields = append(coreGroup.Fields, &field)
	}
	fieldGroups = append(fieldGroups, coreGroup)

	//get tag fields
	tagFields, err := s.getTagFields(ids)
	if err != nil {
		return nil, err
	}
	tagsGroup := model.FieldGroup{
		Name:   model.FieldGroupTags,
		Fields: tagFields,
	}
	fieldGroups = append(fieldGroups, tagsGroup)

	return fieldGroups.AddNullValues(), nil
}

func (s *SQLiteStore) GetResources(ctx context.Context, jsonQuery []byte) (model.ResourcesResponse, error) {
	ids, totalCount, err := s.indexer.findResourceIds(*s.db, s.logger, jsonQuery, true)
	if err != nil {
		return model.ResourcesResponse{}, err
	}
	resources, err := s.getResourcesById(ctx, ids)
	if err != nil {
		return model.ResourcesResponse{}, err
	}
	//update field count to match current query
	allIds := ids
	if totalCount > len(ids) {
		// the response is paginated, but we need all ids to show correct count
		allIds, _, err = s.indexer.findResourceIds(*s.db, s.logger, jsonQuery, false)
		if err != nil {
			return model.ResourcesResponse{}, err
		}
	}
	fields, err := s.getFields(ctx, allIds)
	if err != nil {
		return model.ResourcesResponse{}, err
	}
	return model.ResourcesResponse{Count: totalCount, Resources: resources, FieldGroups: fields}, nil
}

func deleteTags(db *gorm.DB, ids []model.ResourceId) error {
	return db.Table("tags").Where("resource_id in ?", ids).Delete(ids).Error
}

func (s *SQLiteStore) deleteResourcesBefore(before time.Time) (int, error) {

	var rowsAffected int64
	err := s.db.Transaction(func(tx *gorm.DB) error {
		//get the resource ids to delete
		var ids []model.ResourceId
		if err := tx.Table("resources").Select("id").Where("updated_at < ?", before).Find(&ids).Error; err != nil {
			return err
		}

		if len(ids) == 0 {
			//nothing to delete
			return nil
		}

		totalCount, err := resourceCount(tx)
		if err != nil {
			return err
		}
		if totalCount == len(ids) {
			s.logger.Sugar().Warnf("deleting resources before last run would delete all resources (count: %v), ignoring delete.", totalCount)
			// the most common case for this would be the engine encountered a global error,
			// and we are now writting the resource event in the DB to record the error.
			// we woud rather keep the previous resources than deleting them all.
			return nil
		}

		//delete all the tags
		if err := deleteTags(tx, ids); err != nil {
			return err
		}

		//delete the resource indexes and purge the unused columns
		if err := s.indexer.deleteResourceIndexes(tx, ids); err != nil {
			return err
		}
		if err := s.indexer.purgeUnusedColumns(tx); err != nil {
			return err
		}

		//delete the resources
		result := tx.Table("resources").Where("id in ?", ids).Delete(ids)
		if result.Error != nil {
			return result.Error
		}
		rowsAffected = result.RowsAffected
		return nil
	})

	s.logger.Sugar().Infow("Deleting resources: ", zap.Int64("rowsAffected", rowsAffected))

	if err != nil {
		return 0, fmt.Errorf("can't delete resources from database: %w", err)
	}
	return int(rowsAffected), nil
}

func (s *SQLiteStore) EngineStatus(ctx context.Context) (model.Event, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var resourceEvents model.Events
	result := s.db.
		Model(&model.Event{}).
		Order("created_at").
		Find(&resourceEvents, model.Event{RunId: s.runId})
	if result.Error != nil {
		return model.Event{}, fmt.Errorf("error while reading event from database %w", result.Error)
	}
	if len(resourceEvents) == 0 {
		//no event found
		return model.Event{}, nil
	}
	engineEvent := resourceEvents[0]
	engineEvent.ChildEvents = resourceEvents[1:]
	return engineEvent, nil
}

func (s *SQLiteStore) WriteEvent(ctx context.Context, event model.Event) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if event.Type == model.EventTypeEngine && event.Status == model.EventStatusFetching {
		s.runId = event.RunId
		s.fetchedAt = time.Now()
	}
	event.RunId = s.runId

	var existingEvent model.Event
	result := s.db.Model(&model.Event{}).
		Find(&existingEvent, model.Event{
			RunId:        s.runId,
			Type:         event.Type,
			ProviderName: event.ProviderName,
			ResourceType: event.ResourceType,
		})
	if result.Error != nil {
		return fmt.Errorf("error occured while fetching event from database %w", result.Error)
	}
	event.Id = existingEvent.Id
	result = s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&event)
	if result.Error != nil {
		return fmt.Errorf("error occured while upserting events into database %w", result.Error)
	}
	//once engine is complete, we delete all the resources that no longer exist
	if event.Type == model.EventTypeEngine && (event.Status == model.EventStatusFailed || event.Status == model.EventStatusSuccess) {
		_, err := s.deleteResourcesBefore(s.fetchedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

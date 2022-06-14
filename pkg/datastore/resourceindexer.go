package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/a8m/rql"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

const (
	//limit number of resource returned
	DefaultLimit       = 25
	LimitMaxValue      = 2000
	resourceIndexTable = "resource_index"
	//the name the dynamic columns
	columnDynamicName = "col_%v"
	//the maximum number of columns allowed
	//this is a limitation for sqlite https://www.sqlite.org/limits.html
	maxColumns = 2000
)

//resourceIndexer is responsible to index the resources in the DB and provides dynamic querying capabilities
type resourceIndexer struct {
	logger *zap.Logger
	parser *rql.Parser
	//fields that are indexed by FieldName
	fieldColumns fieldColumns
}

//fieldColumn maps a field name to a column name in the query table
type fieldColumn struct {
	//Generated column name, ex: tag1, tag2 ...
	ColumnName string `gorm:"primaryKey"`
	//The actual field name: ex: aws:ec2:fleet-id
	FieldName string
}

//key is field name
type fieldColumns map[string]fieldColumn

func newResourceIndexer(ctx context.Context, logger *zap.Logger, db *gorm.DB) (resourceIndexer, error) {
	qb := resourceIndexer{}
	qb.logger = logger
	qb.fieldColumns = make(fieldColumns)
	err := qb.rebuildDataModel(db)
	return qb, err
}

//newFieldColumns creates the map from a slice of fieldColumn
func newFieldColumns(fields []fieldColumn) fieldColumns {
	result := make(map[string]fieldColumn)
	for _, field := range fields {
		result[field.FieldName] = field
	}
	return result
}

//an explicit field means that the column will be named after the field
func (f fieldColumns) addExplicitFields(names ...string) {
	for _, name := range names {
		f[name] = fieldColumn{
			ColumnName: name,
			FieldName:  name,
		}
	}
}

//a dynamic field would have a generated column name
func (f fieldColumns) addDynamicFields(names ...string) bool {
	newField := false
	for _, name := range names {
		if _, found := f[name]; !found {
			columnName := f.newColumnName()
			f[name] = fieldColumn{
				ColumnName: columnName,
				FieldName:  name,
			}
			newField = true
		}
	}
	return newField
}

//note: this code is inneficient but only done when a new tag key is found (unfrequent)
func (f fieldColumns) containsColumnName(name string) bool {
	for _, v := range f {
		if v.ColumnName == name {
			return true
		}
	}
	return false
}

func (f fieldColumns) newColumnName() string {
	//note: this code is inneficient but only done when a new tag key is found (unfrequent)
	//return a column name that is not used
	var name string
	for i := 1; i == 1 || f.containsColumnName(name); i, name = i+1, fmt.Sprintf(columnDynamicName, i) {
	}
	return name
}

func (f fieldColumns) asSlice() []fieldColumn {
	result := make([]fieldColumn, 0, len(f))
	for _, v := range f {
		result = append(result, v)
	}
	return result
}

//update a query field to map the datamodel
//ex: "aws:ec2:fleet-id" -> "col_3"
func (ri *resourceIndexer) toColumnName(fieldName string) string {
	if fieldCol, ok := ri.fieldColumns[fieldName]; ok {
		return fieldCol.ColumnName
	}
	//do not do any change - rql would throw an unknown field error
	return fieldName
}

//rebuildDataModel updates the query parser and the table schema to includes all the field names
func (ri *resourceIndexer) rebuildDataModel(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("no DB provided")
	}
	if len(ri.fieldColumns) >= maxColumns {
		return fmt.Errorf("maximum number of colums reached: %d", maxColumns)
	}

	if len(ri.fieldColumns) == 0 {
		//first call
		var fieldColumns []fieldColumn
		tables, _ := db.Migrator().GetTables()
		if slices.Contains(tables, "field_columns") {
			//load existing fields from the DB
			db.Find(&fieldColumns)
		}
		ri.fieldColumns = newFieldColumns(fieldColumns)
		//always include these
		ri.fieldColumns.addExplicitFields("id", "type", "region")
	}

	builder := dynamicstruct.NewStruct()
	for _, fieldCol := range ri.fieldColumns {
		structFieldName := cases.Title(language.AmericanEnglish).String(fieldCol.ColumnName)
		if structFieldName == "Id" {
			builder = builder.AddField(structFieldName, "", `gorm:"primaryKey" rql:"filter,sort"`)
		} else {
			builder = builder.AddField(structFieldName, "", `rql:"filter,sort"`)
		}
	}
	//update builder and parse
	resourceIndex := builder.Build().New()
	ri.parser = rql.MustNewParser(rql.Config{
		Model:         resourceIndex,
		FieldSep:      ".",
		DefaultLimit:  DefaultLimit,
		LimitMaxValue: LimitMaxValue,
	})

	//migrate the data to have new columns
	if err := db.Table(resourceIndexTable).AutoMigrate(resourceIndex); err != nil {
		return err
	}

	//save the field column names
	if err := db.AutoMigrate(&fieldColumn{}); err != nil {
		return err
	}
	//update all fields columns
	if err := db.Exec("DELETE FROM field_columns").Error; err != nil {
		return err
	}
	if db := db.Create(ri.fieldColumns.asSlice()); db.Error != nil {
		return db.Error
	}
	return nil
}

func (ri *resourceIndexer) updateQueryFields(jsonQuery []byte) ([]byte, error) {
	var query map[string]interface{}
	err := json.Unmarshal(jsonQuery, &query)
	if err != nil {
		return nil, err
	}
	//update filter
	if obj, ok := query["filter"]; ok {
		if filter, ok := obj.(map[string]interface{}); ok {
			query["filter"] = updateQueryFilter(filter, ri.toColumnName)
		}
	}
	//update sort
	if obj, ok := query["sort"]; ok {
		if sort, ok := obj.([]interface{}); ok {
			query["sort"] = updateQuerySort(sort, ri.toColumnName)
		}
	}
	jsonQuery, err = json.Marshal(query)
	if err != nil {
		return nil, err
	}
	return jsonQuery, nil
}

func updateQueryFilter(filter map[string]interface{}, f func(string) string) map[string]interface{} {
	/*{
	  "filter":{
	    "type":"ec2.Volume",
	    "$or": [
	      { "team": "marketplace" },
	      { "team": "shipping" }
	    ]
	  }
	}
	->will update the fields: type & team
	*/
	result := make(map[string]interface{}, len(filter))
	for k, v := range filter {
		if k == "$or" || k == "$and" {
			if slice, ok := v.([]interface{}); ok {
				for i, val := range slice {
					if _map, ok := val.(map[string]interface{}); ok {
						slice[i] = updateQueryFilter(_map, f)
					}
				}
			}
		} else {
			//update key
			k = f(k)
		}
		//update the result
		result[k] = v
	}
	return result
}
func updateQuerySort(sort []interface{}, f func(string) string) []interface{} {
	/*
		{
			"sort": ["-region"]
		}
		-> will update the field: region
	*/
	result := make([]interface{}, len(sort))
	for i, val := range sort {
		//update the values for sort
		if s, ok := val.(string); ok && len(s) > 0 {
			//special logic for sort to handle '+' or '-' prefix
			if strings.ContainsRune("+-", rune(s[0])) {
				result[i] = s[0:1] + f(string(s[1:]))
			} else {
				result[i] = f(string(s))
			}
		}
	}
	return result
}

func (ri *resourceIndexer) deleteResourceIndexes(db *gorm.DB, ids []model.ResourceId) error {
	err := db.Table(resourceIndexTable).Where("id in ?", ids).Delete(ids).Error
	if err != nil {
		return fmt.Errorf("could not delete the resources indexes: %w", err)
	}
	return nil
}

//purgeUnusedColumns delete the unused columns so they can be used for other tag keys
func (ri *resourceIndexer) purgeUnusedColumns(db *gorm.DB) error {
	purged := false
	for key, col := range ri.fieldColumns {
		if !strings.HasPrefix(col.ColumnName, fmt.Sprintf(columnDynamicName, "")) {
			//only remove dynamic columns
			continue
		}
		//count the values defined for the current column
		var count int64
		err := db.Table(resourceIndexTable).Where(fmt.Sprintf("%v is not null", col.ColumnName)).Count(&count).Error
		if err != nil {
			return err
		}
		if count == 0 {
			//this column has no value set, it can be removed
			ri.logger.Sugar().Debugf("the column %v (tag: %v) will be removed from %v", col.ColumnName, col.FieldName, resourceIndexTable)
			delete(ri.fieldColumns, key)
			purged = true
		}
	}

	if purged {
		//only if something was changed we would rebuild the model
		return ri.rebuildDataModel(db)
	}
	return nil
}

//writeResourceIndexes will insert a new resource_index row for each Resource
func (ri *resourceIndexer) writeResourceIndexes(ctx context.Context, db *gorm.DB, resources []*model.Resource) error {

	//get all the possible tag keys - once per batch
	rebuildModel := false
	for _, r := range resources {
		for _, tag := range r.Tags {
			if ri.fieldColumns.addDynamicFields(tag.Key) {
				//new tag key found - need to rebuild the model
				rebuildModel = true
			}
		}
	}
	if rebuildModel {
		if err := ri.rebuildDataModel(db); err != nil {
			return err
		}
	}

	//create the rows to insert in memory
	rows := make([]map[string]interface{}, len(resources))
	var ids []model.ResourceId
	for i, r := range resources {
		row := make(map[string]interface{})
		ids = append(ids, model.ResourceId(r.Id))
		row["id"] = r.Id
		row["type"] = r.Type
		row["region"] = r.Region
		//add the tags
		for _, tag := range r.Tags {
			//ex: row["Col23"]="Jordan"
			row[ri.fieldColumns[tag.Key].ColumnName] = tag.Value
		}
		rows[i] = row
	}

	//delete all previous resource_index if they exist
	if err := ri.deleteResourceIndexes(db, ids); err != nil {
		return err
	}

	//create the resource_index
	if err := db.Table(resourceIndexTable).Create(rows).Error; err != nil {
		return fmt.Errorf("could not create the resource indexes: %w", err)
	}

	return nil
}

func (ri *resourceIndexer) parse(jsonQuery []byte) (*rql.Params, error) {
	ri.logger.Sugar().Debugw("received",
		zap.String("query", string(jsonQuery)),
	)
	// update query field names to map to the data model
	jsonQuery, err := ri.updateQueryFields(jsonQuery)
	if err != nil {
		return nil, err
	}
	ri.logger.Sugar().Debugw("updated",
		zap.String("query", string(jsonQuery)),
	)
	params, err := ri.parser.Parse(jsonQuery)
	if err != nil {
		return nil, err
	}
	//hand null values
	return replaceNullValues(params), nil
}

func replaceNullValues(p *rql.Params) *rql.Params {
	for i, arg := range p.FilterArgs {
		if s, ok := arg.(string); ok {
			if s == model.NullValue {
				p.FilterExp = replaceWith(p.FilterExp, "=", "is", "?", i)
				p.FilterArgs[i] = nil
			}
			if s == model.NotNullValue {
				p.FilterExp = replaceWith(p.FilterExp, "=", "is not", "?", i)
				p.FilterArgs[i] = nil
			}
		}
	}
	return p
}

func replaceWith(s string, old string, new string, sep string, i int) string {
	parts := strings.Split(s, sep)
	parts[i] = strings.Replace(parts[i], old, new, 1)
	return strings.Join(parts, sep)
}

//findResourceIds finds resources using a RQL query
//see https://github.com/a8m/rql#getting-started for the syntax
func (ri *resourceIndexer) findResourceIds(db gorm.DB, logger *zap.Logger, jsonQuery []byte) ([]model.ResourceId, int, error) {
	if len(jsonQuery) == 0 {
		//use en empty json if nothing is set - this will use the default limit
		jsonQuery = []byte(`{}`)
	}
	p, err := ri.parse(jsonQuery)
	if err != nil {
		return nil, 0, err
	}

	var resourceIds []model.ResourceId
	result := db.Table(resourceIndexTable).
		Select("id").
		Where(p.FilterExp, p.FilterArgs...).
		Offset(p.Offset).
		Limit(p.Limit).
		Order(p.Sort).
		Find(&resourceIds)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	var count int64
	resultCount := db.Table(resourceIndexTable).
		Select("id").
		Where(p.FilterExp, p.FilterArgs...).
		Count(&count)

	if resultCount.Error != nil {
		return nil, 0, result.Error
	}
	return resourceIds, int(count), nil
}

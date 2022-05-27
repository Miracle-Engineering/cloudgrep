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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	//limit number of resource returned
	defaultLimit       = 25
	limitMaxValue      = 2000
	resourceIndexTable = "resource_index"
	//the name the dynamic columns
	columnDynamicName = "col_%d"
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

type fieldColumns map[string]fieldColumn

func newResourceIndexer(ctx context.Context, logger *zap.Logger, db *gorm.DB) (resourceIndexer, error) {
	qb := resourceIndexer{}
	qb.logger = logger
	qb.fieldColumns = make(fieldColumns)
	err := qb.rebuildDataModel(ctx, db)
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
func (qb *resourceIndexer) toColumnName(fieldName string) string {
	if fieldCol, ok := qb.fieldColumns[fieldName]; ok {
		return fieldCol.ColumnName
	}
	//do not do any change - rql would throw an unknown field error
	return fieldName
}

//rebuildDataModel updates the query parser and the table schema to includes all the field names
func (qb *resourceIndexer) rebuildDataModel(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("no DB provided")
	}
	if len(qb.fieldColumns) >= maxColumns {
		return fmt.Errorf("maximum number of colums reached: %d", maxColumns)
	}

	if len(qb.fieldColumns) == 0 {
		//first call
		var fieldColumns []fieldColumn
		//load existing fields from the DB
		db.Find(&fieldColumns)
		qb.fieldColumns = newFieldColumns(fieldColumns)
		//always include these
		qb.fieldColumns.addExplicitFields("id", "type", "region")
	}

	builder := dynamicstruct.NewStruct()
	for _, fieldCol := range qb.fieldColumns {
		structFieldName := cases.Title(language.AmericanEnglish).String(fieldCol.ColumnName)
		if structFieldName == "Id" {
			builder = builder.AddField(structFieldName, "", `gorm:"primaryKey" rql:"filter,sort"`)
		} else {
			builder = builder.AddField(structFieldName, "", `rql:"filter,sort"`)
		}
	}
	//update builder and parse
	resourceIndex := builder.Build().New()
	qb.parser = rql.MustNewParser(rql.Config{
		Model:         resourceIndex,
		FieldSep:      ".",
		DefaultLimit:  defaultLimit,
		LimitMaxValue: limitMaxValue,
	})

	//migrate the data to have new columns
	if err := db.Table(resourceIndexTable).AutoMigrate(resourceIndex); err != nil {
		return err
	}

	//save the field column names
	//TODO we should purge this table to remove the unused fields (resources no longer exist)
	//we can address this when we support Refresh/update resources
	if err := db.AutoMigrate(&fieldColumn{}); err != nil {
		return err
	}
	if db := db.Clauses(clause.OnConflict{DoNothing: true}).Create(qb.fieldColumns.asSlice()); db.Error != nil {
		return db.Error
	}
	return nil
}

func (qb *resourceIndexer) updateQueryFields(jsonQuery []byte) ([]byte, error) {
	var query map[string]interface{}
	err := json.Unmarshal(jsonQuery, &query)
	if err != nil {
		return nil, err
	}
	//update filter
	if obj, ok := query["filter"]; ok {
		if filter, ok := obj.(map[string]interface{}); ok {
			query["filter"] = updateQueryFilter(filter, qb.toColumnName)
		}
	}
	//update sort
	if obj, ok := query["sort"]; ok {
		if sort, ok := obj.([]interface{}); ok {
			query["sort"] = updateQuerySort(sort, qb.toColumnName)
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
		if k == "$or" {
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
		//update the values
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

//writeResourceIndexes will insert a new resource_index row for each Resource
func (qb *resourceIndexer) writeResourceIndexes(ctx context.Context, db *gorm.DB, resources []*model.Resource) error {

	//get all the possible tag keys - once per batch
	rebuildModel := false
	for _, r := range resources {
		for _, tag := range r.Tags {
			if qb.fieldColumns.addDynamicFields(tag.Key) {
				//new tag key found - need to rebuild the model
				rebuildModel = true
			}
		}
	}
	if rebuildModel {
		if err := qb.rebuildDataModel(ctx, db); err != nil {
			return err
		}
	}

	//write all the resource_index rows
	rows := make([]map[string]interface{}, len(resources))
	for i, r := range resources {
		row := make(map[string]interface{})
		row["id"] = r.Id
		row["type"] = r.Type
		row["region"] = r.Region
		//add the tags
		for _, tag := range r.Tags {
			//ex: row["Col23"]="Jordan"
			row[qb.fieldColumns[tag.Key].ColumnName] = tag.Value
		}
		rows[i] = row
	}
	result := db.Table(resourceIndexTable).Create(rows)
	if result.Error != nil {
		return fmt.Errorf("could not index the resources: %w", result.Error)
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
			if s == "[null]" {
				p.FilterExp = replaceWith(p.FilterExp, "=", "is", "?", i)
				p.FilterArgs[i] = nil
			}
			if s == "[not null]" {
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
func (ri *resourceIndexer) findResourceIds(db gorm.DB, logger *zap.Logger, jsonQuery []byte) ([]resourceId, error) {
	if len(jsonQuery) == 0 {
		//use en empty json if nothing is set - this will use the default limit
		jsonQuery = []byte(`{}`)
	}
	p, err := ri.parse(jsonQuery)
	if err != nil {
		return nil, err
	}

	var resourceIds []resourceId
	result := db.Table(resourceIndexTable).
		Select("id").
		Where(p.FilterExp, p.FilterArgs...).
		Offset(p.Offset).
		Limit(p.Limit).
		Order(p.Sort).
		Find(&resourceIds)

	if result.Error != nil {
		return nil, result.Error
	}
	return resourceIds, nil
}

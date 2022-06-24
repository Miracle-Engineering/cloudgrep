package model

import (
	"fmt"
	"strconv"
)

const (
	//NullValue used in a query, means that the resource should not have this field defined
	NullValue = "(null)"
	//NullValue used in a query, means that the resource should have this field defined
	NotNullValue = "(not null)"

	//CountValueIgnored means that the current value is ignored in the current query
	CountValueIgnored = "-"

	FieldGroupCore = "core"
	FieldGroupTags = "tags"
)

//FieldGroup regroups some fields. Ex: "Tags"
type FieldGroup struct {
	Name   string `json:"name"`
	Fields Fields `json:"fields"`
}
type FieldGroups []FieldGroup

//Field is a searchable attribute on a resource.
//It also includes the possible values and their respective count.
type Field struct {
	Name   string      `json:"name"`
	Count  int         `json:"count"`
	Values FieldValues `json:"values"`
}

//FieldValue is a value associated with a field.
//The count is the number of resources with this field value.
type FieldValue struct {
	Value string `json:"value"`
	Count string `json:"count"`
}

type Fields []*Field
type FieldValues []*FieldValue

func (fgs FieldGroups) FindGroup(group string) *FieldGroup {
	for _, fg := range fgs {
		if fg.Name == group {
			return &fg
		}
	}
	return nil
}

func (fgs FieldGroups) FindField(group string, name string) *Field {
	fg := fgs.FindGroup(group)
	if fg != nil {
		for _, field := range fg.Fields {
			if field.Name == name {
				return field
			}
		}
	}
	return nil
}

func (fv *FieldValues) Find(value string) *FieldValue {
	for _, v := range *fv {
		if v.Value == value {
			return v
		}
	}
	return nil
}

func (fv FieldValues) Count() int {
	count := 0
	for _, v := range fv {
		if v.Count != CountValueIgnored {
			c, _ := strconv.Atoi(v.Count)
			count = count + c
		}
	}
	return count
}

//count returns the number of resources
func (fgs FieldGroups) count() int {
	//assuming each resource always has type set
	if field := fgs.FindField("core", "type"); field != nil {
		return field.Count
	}
	return 0
}

//AddNullValues adds the (null) value for each Field, it is used by the API to allow filtering on resources without the field.
// If a field is always defined (ex: type), do not include the (null) value as it would mean excluding all resources from a query.
func (fgs FieldGroups) AddNullValues() FieldGroups {
	var result []FieldGroup
	totalCount := fgs.count()
	for _, group := range fgs {
		fields := make([]*Field, 0)
		for _, field := range group.Fields {

			nullCount := totalCount - field.Count
			//do not show null if all resource would be excluded
			if nullCount > 0 {
				field.Values = append(field.Values,
					&FieldValue{
						Value: NullValue,
						Count: fmt.Sprint(nullCount),
					})
			}
			fields = append(fields, field)
		}
		result = append(result, FieldGroup{Name: group.Name, Fields: fields})
	}
	return result
}

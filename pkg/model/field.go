package model

type FieldGroup struct {
	Name   string `json:"name"`
	Fields Fields `json:"fields"`
}
type FieldGroups []FieldGroup

type Field struct {
	Name   string      `json:"name"`
	Count  int         `json:"count"`
	Values FieldValues `json:"values"`
}

type FieldValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type Fields []Field
type FieldValues []FieldValue

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
				return &field
			}
		}
	}
	return nil
}

package model

type Field struct {
	Name   string      `json:"name"`
	Group  string      `json:"group"`
	Count  int         `json:"count"`
	Values FieldValues `json:"values"`
}

type FieldValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type Fields []Field
type FieldValues []FieldValue

func (fields Fields) Find(name string) *Field {
	for _, f := range fields {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

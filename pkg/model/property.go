package model

import "go.uber.org/zap/zapcore"

type Property struct {
	ResourceId string `json:"-" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"primaryKey"`
	Value      string `json:"value"`
}
type Properties []Property

func (p Property) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", p.Name)
	enc.AddString("value", p.Value)
	return nil
}

func (p Properties) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, prop := range p {
		if err := enc.AppendObject(prop); err != nil {
			return err
		}
	}
	return nil
}

//clean removes unexported fields
func (p Property) clean() Property {
	return Property{
		Name:  p.Name,
		Value: p.Value,
	}
}
func (p Properties) Clean() Properties {
	var props Properties
	for _, prop := range p {
		props = append(props, prop.clean())
	}
	return props
}

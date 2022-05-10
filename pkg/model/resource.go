package model

import "go.uber.org/zap/zapcore"

//TODO store provider info in resource (needed when we can have more than one provider)
type Resource struct {
	Id         string     `json:"id" gorm:"primaryKey"`
	Region     string     `json:"region"`
	Type       string     `json:"type"`
	Tags       Tags       `json:"tags"`
	Properties Properties `json:"properties"`
}

func (r Resource) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", r.Id)
	enc.AddString("region", r.Region)
	enc.AddString("type", r.Type)
	//do not display tags and regions by default - too verbose
	return nil
}

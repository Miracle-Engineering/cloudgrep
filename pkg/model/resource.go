package model

import (
	"time"

	"go.uber.org/zap/zapcore"
	"gorm.io/datatypes"
)

//TODO store provider info in resource (needed when we can have more than one provider)
type Resource struct {
	Id        string         `json:"id" gorm:"primaryKey"`
	DisplayId string         `json:"displayId"`
	AccountId string         `json:"accountId"`
	Region    string         `json:"region"`
	Type      string         `json:"type"`
	Tags      Tags           `json:"tags"`
	RawData   datatypes.JSON `json:"rawData"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type Resources []*Resource
type ResourcesResponse struct {
	Count       int         `json:"count"`
	FieldGroups FieldGroups `json:"fieldGroups"`
	Resources   Resources   `json:"resources"`
}
type ResourceId string

func (r Resource) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", r.Id)
	enc.AddString("region", r.Region)
	enc.AddString("type", r.Type)
	//do not display tags and regions and raw data by default - too verbose
	return nil
}

func (r *Resource) clean() *Resource {
	//replace with the default value
	r.UpdatedAt = Resource{}.UpdatedAt
	return r
}

//Clean removes the generated fields - useful for testing where the same instance is reused
func (rs Resources) Clean() Resources {
	var result Resources
	for _, r := range rs {
		result = append(result, r.clean())
	}
	return result
}

//FindById finds a resource by ID, return nil if not found
func (rs Resources) FindById(id string) *Resource {
	for _, r := range rs {
		if r.Id == id {
			return r
		}
	}
	return nil
}

func (rs Resources) Ids() []ResourceId {
	return ResourceIds(rs)
}

func ResourceIds(resources []*Resource) []ResourceId {
	ids := make([]ResourceId, len(resources))
	for i, r := range resources {
		ids[i] = ResourceId(r.Id)
	}
	return ids
}

package model

//TODO store provider info in resource (needed when we can have more than one provider)
type Resource struct {
	Id         string     `json:"id" gorm:"primaryKey"`
	Region     string     `json:"region"`
	Type       string     `json:"type"`
	Tags       []Tag      `json:"tags"`
	Properties []Property `json:"properties"`
}

type Tag struct {
	ResourceId string `gorm:"primaryKey"`
	Key        string `json:"key" gorm:"primaryKey"`
	Value      string `json:"value"`
}

type Property struct {
	ResourceId string `gorm:"primaryKey"`
	Name       string `json:"name" gorm:"primaryKey"`
	Value      string `json:"value"`
}

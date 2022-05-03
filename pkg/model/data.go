package model

type Resource struct {
	Id         string     `json:"id"`
	Region     string     `json:"region"`
	Type       string     `json:"type"`
	Tags       []Tag      `json:"tags"`
	Properties []Property `json:"properties"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

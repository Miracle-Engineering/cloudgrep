package config

type Service struct {
	Name           string
	ServicePackage string `yaml:"servicePackage"`
	Global         bool   `yaml:"global"`
	Types          []Type `yaml:"types"`
}

type Type struct {
	Name        string    `yaml:"name"`
	ListAPI     ListAPI   `yaml:"listApi"`
	GetTagsAPI  GetTagAPI `yaml:"getTagsApi"`
	Description string    `yaml:"description"`
	Global      *bool     `yaml:"global"`
}

type TagField struct {
	Style   string      `yaml:"style"`
	Pointer bool        `yaml:"pointer"`
	Field   NestedField `yaml:"field"`
	Key     string      `yaml:"key"`
	Value   string      `yaml:"value"`
}

type ListAPI struct {
	Call       string    `yaml:"call"`
	Pagination bool      `yaml:"pagination"`
	OutputKey  []string  `yaml:"outputKey"`
	IDField    Field     `yaml:"id"`
	Tags       *TagField `yaml:"tags"`
}

type GetTagAPI struct {
	ResourceType         string      `yaml:"type"`
	Call                 string      `yaml:"call"`
	InputIDField         Field       `yaml:"inputIDField"`
	OutputKey            NestedField `yaml:"outputKey"`
	TagField             *TagField   `yaml:"tags"`
	AllowedAPIErrorCodes []string    `yaml:"allowedApiErrorCodes"`
}

func (c GetTagAPI) Has() bool {
	return c.Call != ""
}

type Root struct {
	Services []string `yaml:"services"`
}

type Config struct {
	Services []Service
}

type NestedField []Field

type Field struct {
	Name      string `yaml:"name"`
	SliceType string `yaml:"sliceType"`
	Pointer   bool   `yaml:"pointer"`
}

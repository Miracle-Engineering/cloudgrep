package config

type ServiceConfig struct {
	Name           string
	ServicePackage string       `yaml:"servicePackage"`
	Types          []TypeConfig `yaml:"types"`
}

type TypeConfig struct {
	Name        string          `yaml:"name"`
	ListAPI     ListAPIConfig   `yaml:"listApi"`
	GetTagsAPI  GetTagAPIConfig `yaml:"getTagsApi"`
	Description string          `yaml:"description"`
	Global      bool            `yaml:"global"`
}

// type TagStyle string

// const (
// 	TagStyleMap   TagStyle = "map"
// 	TagStyleField TagStyle = "field"
// )

type TagField struct {
	Style   string      `yaml:"style"`
	Pointer bool        `yaml:"pointer"`
	Field   NestedField `yaml:"field"`
	Key     string      `yaml:"key"`
	Value   string      `yaml:"value"`
}

func (t TagField) Zero() bool {
	return t.Field.Empty()
}

type ListAPIConfig struct {
	Call       string   `yaml:"call"`
	Pagination bool     `yaml:"pagination"`
	OutputKey  []string `yaml:"outputKey"`
	IDField    Field    `yaml:"id"`
	Tags       TagField `yaml:"tags"`
}

type GetTagAPIConfig struct {
	ResourceType         string      `yaml:"type"`
	Call                 string      `yaml:"call"`
	InputIDField         Field       `yaml:"inputIDField"`
	OutputKey            NestedField `yaml:"outputKey"`
	TagField             TagField    `yaml:"tags"`
	AllowedAPIErrorCodes []string    `yaml:"allowedApiErrorCodes"`
}

func (c GetTagAPIConfig) Has() bool {
	return c.Call != ""
}

type RootConfig struct {
	Services []string `yaml:"services"`
}

type Config struct {
	Services []ServiceConfig
}

func newConfig() *Config {
	return &Config{
		// Services: map[string]ServiceConfig{},
	}
}

type NestedField []Field

func (f NestedField) Last() Field {
	if len(f) == 0 {
		return Field{}
	}

	return f[len(f)-1]
}

type Field struct {
	Name      string `yaml:"name"`
	SliceType string `yaml:"sliceType"`
	Pointer   bool   `yaml:"pointer"`
}

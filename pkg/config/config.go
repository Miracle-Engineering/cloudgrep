package config

import (
	_ "embed"
	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var EmbedConfig []byte

type Config struct {
	Providers []Provider `yaml:"providers"`
	Datastore Datastore  `yaml:"datastore"`
	Web       Web        `yaml:"web"`
	// Adding regions as where cli regions override will be stored
	Regions []string
}

type Provider struct {
	Cloud   string   `yaml:"cloud"`
	Regions []string `yaml:"regions"`
}

type Datastore struct {
	Type           string `yaml:"type"`
	SkipRefresh    bool   `yaml:"skipRefresh"`
	DataSourceName string `yaml:"dataSourceName"`
}

type Web struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Prefix   string `yaml:"prefix"`
	SkipOpen bool   `yaml:"skip_open"`
}

func GetDefault() (Config, error) {
	var err error
	var config Config

	err = yaml.Unmarshal(EmbedConfig, &config)
	if err != nil {
		return Config{}, err
	}
	config.Regions = []string{}
	return config, err
}

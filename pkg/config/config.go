package config

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
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

func (p *Provider) String() string {
	if strings.Join(p.Regions, "-") == "" {
		return p.Cloud
	}
	return fmt.Sprintf("%s-%s", p.Cloud, strings.Join(p.Regions, "-"))
}

type Datastore struct {
	Type           string `yaml:"type"`
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

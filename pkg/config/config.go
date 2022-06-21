package config

import (
	_ "embed"
	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var EmbedConfig []byte

// Config represents all the user-configurable settings for cloudgrep. One such structure is loaded at runtime and
// is populated through the cli arguments, user-provider config file, or a preset default, with values resolved
// in that order of precedence. To see the default, please refer to
// https://github.com/run-x/cloudgrep/blob/main/pkg/config/config.yaml
type Config struct {
	// Providers represents the providers to be scanned by cloudgrep
	Providers []Provider `yaml:"providers"`
	// Datastore is the datastore to be used
	Datastore Datastore `yaml:"datastore"`
	// Web is the web specs to be used
	Web Web `yaml:"web"`
	// Adding regions as where cli regions override will be stored
	Regions []string
}

// Provider represents a cloud provider cloudgrep will scan w/ the current credentials
type Provider struct {
	// Cloud is the type of the cloud provider (currently only AWS is supported)
	Cloud string `yaml:"cloud"`
	// Regions is the list of different regions within the cloud provider to scan
	Regions []string `yaml:"regions"`
}

// Datastore represents the specs cloudgrep uses for creating and/or connecting to the datastore/database used.
type Datastore struct {
	// Type is the kind of datastore to be used by cloudgrep (currently only supports SQLite)
	Type string `yaml:"type"`
	// SkipRefresh determines whether to refresh the data (i.e. scan the cloud) on startup.
	SkipRefresh bool `yaml:"skipRefresh"`
	// DataSourceName is the Type-specific data source name or uri for connecting to the desired data source
	DataSourceName string `yaml:"dataSourceName"`
}

// Web represents the specs cloudgrep uses for creating the webapp server
type Web struct {
	// Host is the host the server is running as
	Host string `yaml:"host"`
	// Port is the port the server is running in
	Port int `yaml:"port"`
	// Prefix is the url prefix the server uses
	Prefix string `yaml:"prefix"`
	// SkipOpen determines whether to automatically open the webui on startup
	SkipOpen bool `yaml:"skipOpen"`
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

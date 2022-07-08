package config

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

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
	// Adding regions as where cli regions override is stored
	Regions []string
	// Adding regions as where cli profiles override is stored
	Profiles []string
	// Adding roles as where cli profiles override is stored
	RoleArns []string
}

// Provider represents a cloud provider cloudgrep will scan w/ the current credentials
type Provider struct {
	// Cloud is the type of the cloud provider (currently only AWS is supported)
	Cloud string `yaml:"cloud"`
	// Regions is the list of different regions within the cloud provider to scan
	Regions []string `yaml:"regions"`
	// Profile is the AWS profile to use, if not set use the default profile
	Profile string `yaml:"profile"`
	// RoleArn is the AWS role ARN to assume. If not set, do not assume any role
	RoleArn string `yaml:"roleArn"`
}

func (p *Provider) String() string {
	if len(p.Regions) == 0 {
		return p.Cloud
	}
	return fmt.Sprintf("%s-%s", p.Cloud, strings.Join(p.Regions, "-"))
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
	return config, nil
}

func ReadFile(file string) (Config, error) {
	cfg, err := GetDefault()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return cfg, err
	}
	if err != nil {
		return cfg, err
	}
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

//Load will load the config files making sure that the options override are correct
func (c *Config) Load() error {
	err := c.loadProfiles()
	if err != nil {
		return err
	}
	err = c.loadRoleArns()
	if err != nil {
		return err
	}
	err = c.loadRegions()
	if err != nil {
		return err
	}
	return err
}

//loadRegions will replace the providers/regions if the --regions option if set
func (c *Config) loadRegions() error {
	if len(c.Regions) == 0 {
		return nil
	}
	providers := make([]Provider, 0)
	for _, provider := range c.Providers {
		provider.Regions = c.Regions
		providers = append(providers, provider)
	}
	c.Providers = providers
	return nil
}

//loadProfiles will replace the providers/profile if the --providers option if set
func (c *Config) loadProfiles() error {
	if len(c.Profiles) == 0 {
		return nil
	}
	//create the provider for each profile
	providers := make([]Provider, 0)
	for _, profile := range c.Profiles {
		for _, provider := range c.Providers {
			if provider.Profile != "" {
				return fmt.Errorf("the config file already defines a profile, using the option `--profiles` is not supported")
			}
			providerNew := provider
			providerNew.Profile = profile
			providers = append(providers, providerNew)
		}
	}
	c.Providers = providers
	return nil
}

//loadRoleArns will replace the providers/roleArn if the --role-arns option if set
func (c *Config) loadRoleArns() error {
	if len(c.RoleArns) == 0 {
		return nil
	}
	//create the provider for each role
	providers := make([]Provider, 0)
	for _, roleArn := range c.RoleArns {
		for _, provider := range c.Providers {
			if provider.RoleArn != "" {
				return fmt.Errorf("the config file already defines a roleArn, using the option `--role-arns` is not supported")
			}
			providerNew := provider
			providerNew.RoleArn = roleArn
			providers = append(providers, providerNew)
		}
	}
	c.Providers = providers
	return nil
}

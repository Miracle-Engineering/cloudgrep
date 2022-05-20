package main

import (
	"fmt"
	"path"

	"gopkg.in/yaml.v2"
)

type ServiceConfig struct {
	Name  string
	Types []TypeConfig `yaml:"types"`
}

type TypeConfig struct {
	Name        string          `yaml:"name"`
	ListAPI     ListAPIConfig   `yaml:"listApi"`
	GetTagsAPI  GetTagAPIConfig `yaml:"getTagsApi"`
	Description string          `yaml:"description"`
}

type ListAPIConfig struct {
	Call       string   `yaml:"call"`
	Pagination bool     `yaml:"pagination"`
	OutputKey  []string `yaml:"outputKey"`
	IDField    string   `yaml:"id"`
	TagField   string   `yaml:"tags"`
}

type GetTagAPIConfig struct {
	Call         string `yaml:"call"`
	InputIDField string `yaml:"inputIDField"`
	OutputKey    string `yaml:"outputKey"`
}

type ServicesConfig struct {
	Services []string `yaml:"services"`
}

func loadServices() map[string]ServiceConfig {
	serviceListPath := path.Join("services", "config.yaml")

	serviceListRaw, err := content.ReadFile(serviceListPath)
	if err != nil {
		panic(fmt.Errorf("cannot read service list: %w", err))
	}

	var config ServicesConfig
	err = yaml.UnmarshalStrict(serviceListRaw, &config)
	if err != nil {
		panic(fmt.Errorf("cannot parse service list: %w", err))
	}

	services := make(map[string]ServiceConfig)
	for _, name := range config.Services {
		services[name] = loadService(name)
	}

	return services
}

func loadService(name string) ServiceConfig {
	servicePath := path.Join("services", name+".yaml")
	serviceRaw, err := content.ReadFile(servicePath)
	if err != nil {
		panic(fmt.Errorf("cannot read service %s: %w", name, err))
	}

	var config ServiceConfig
	err = yaml.UnmarshalStrict(serviceRaw, &config)
	if err != nil {
		panic(fmt.Errorf("cannot parse service %s: %w", name, err))
	}

	config.Name = name

	return config
}

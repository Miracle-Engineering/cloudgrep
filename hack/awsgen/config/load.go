package config

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func Load(rootConfigPath string) (*Config, error) {
	root, err := loadRoot(rootConfigPath)
	if err != nil {
		return nil, err
	}

	rootDir := path.Dir(rootConfigPath)

	config := newConfig()

	for _, serviceName := range root.Services {
		service, err := loadService(rootDir, serviceName)
		if err != nil {
			return nil, fmt.Errorf("cannot load service %s: %w", serviceName, err)
		}

		config.Services = append(config.Services, service)
	}

	return config, nil
}

func loadRoot(configPath string) (RootConfig, error) {
	c := RootConfig{}

	err := loadConfigYaml(configPath, &c)
	return c, err
}

func loadService(dir string, name string) (ServiceConfig, error) {
	c := ServiceConfig{
		Name:           name,
		ServicePackage: name,
	}

	configPath := path.Join(dir, name+".yaml")

	err := loadConfigYaml(configPath, &c)

	return c, err
}

func loadConfigYaml(path string, val any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot load config at %s: %w", path, err)
	}

	d := yaml.NewDecoder(f)
	d.KnownFields(true)

	err = d.Decode(val)
	if err != nil {
		return fmt.Errorf("cannot parse config at %s: %w", path, err)
	}

	return nil
}

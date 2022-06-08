package config

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func Load(rootConfigPath string) (*Config, error) {
	root, err := loadRoot(rootConfigPath)
	if err != nil {
		return nil, err
	}

	rootDir := path.Dir(rootConfigPath)

	config := &Config{}

	sort.Strings(root.Services)

	for _, serviceName := range root.Services {
		service, err := loadService(rootDir, serviceName)
		if err != nil {
			return nil, fmt.Errorf("cannot load service %s: %w", serviceName, err)
		}

		config.Services = append(config.Services, service)
	}

	return config, nil
}

func loadRoot(configPath string) (Root, error) {
	c := Root{}

	err := loadConfigYaml(configPath, &c)
	return c, err
}

func loadService(dir string, name string) (Service, error) {
	c := Service{
		Name:           name,
		ServicePackage: name,
	}

	configPath := path.Join(dir, name+".yaml")

	err := loadConfigYaml(configPath, &c)
	if err != nil {
		return c, err
	}

	sort.Slice(c.Types, func(i, j int) bool {
		return strings.Compare(c.Types[i].Name, c.Types[j].Name) < 0
	})

	return c, err
}

func loadConfigYaml(path string, val any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot load config at %s: %w", path, err)
	}

	defer f.Close()

	d := yaml.NewDecoder(f)
	d.KnownFields(true)

	err = d.Decode(val)
	if err != nil {
		return fmt.Errorf("cannot parse config at %s: %w", path, err)
	}

	return nil
}

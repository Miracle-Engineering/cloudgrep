package demo

import (
	_ "embed"
	"github.com/run-x/cloudgrep/pkg/config"
	"gopkg.in/yaml.v2"
)

//go:embed demo.yaml
var DemoConfig []byte

//go:embed demo.db
var DemoDB []byte

func GetDemoConfig() (config.Config, error) {
	cfg, err := config.GetDefault()
	if err != nil {
		return config.Config{}, err
	}

	err = yaml.Unmarshal(DemoConfig, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, err
}

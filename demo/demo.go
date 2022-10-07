package demo

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/juandiegopalomino/cloudgrep/pkg/config"
	"gopkg.in/yaml.v2"
)

//go:embed demo.yaml
var DemoConfig []byte

//go:embed demo.db
var DemoDB []byte

// GetDemoConfig returns the demo config - this will copy over the demo file to a temporary file
// you can cleanup this file by deleting cfg.Datastore.DataSourceName when done
func GetDemoConfig() (config.Config, error) {
	cfg, err := config.GetDefault()
	if err != nil {
		return config.Config{}, err
	}

	err = yaml.Unmarshal(DemoConfig, &cfg)
	if err != nil {
		return config.Config{}, err
	}

	//when running the demo from the binary the embeded file is copied over to a local temp dir
	file, err := os.CreateTemp("", "cloudgrepdemodb")
	if err != nil {
		return config.Config{}, err
	}
	_, err = file.Write(DemoDB)
	if err != nil {
		return config.Config{}, err
	}
	cfg.Datastore.DataSourceName, err = filepath.Abs(file.Name())
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

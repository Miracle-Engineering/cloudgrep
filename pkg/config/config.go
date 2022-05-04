package config

import (
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"

	"github.com/run-x/cloudgrep/pkg/options"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var embedConfig []byte

type Config struct {
	Providers []Provider `yaml:"providers"`
	Datastore Datastore  `yaml:"datastore"`
	Web       Web        `yaml:"web"`
	Logging   Logging    `yaml:"logging"`
}

type Provider struct {
	Cloud string `yaml:"cloud"`
}

type Datastore struct {
	Type string `yaml:"type"`
}

type Web struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Prefix string `yaml:"prefix"`
}

const loggingModeDev = "dev"

const loggingModeProd = "prod"

type Logging struct {
	Mode   string `yaml:"mode"`
	Logger *zap.Logger
}

//New creates a new configuration
func New(ctx context.Context, options options.Options) (Config, error) {
	var data []byte
	var err error

	if options.Config == "" {
		//Read the default configuration
		data = embedConfig
	} else {
		data, err = ioutil.ReadFile(options.Config)
	}

	if err != nil {
		return Config{}, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	// some config can be overwritten by setting some CLI options
	if options.HTTPHost != "" {
		config.Web.Host = options.HTTPHost
	}
	if options.HTTPPort > 0 {
		config.Web.Port = int(options.HTTPPort)
	}
	if options.Prefix != "" {
		config.Web.Prefix = options.Prefix
	}

	//init logger
	var logger *zap.Logger
	switch config.Logging.Mode {
	case loggingModeDev:
		logger, err = zap.NewDevelopment()
	case loggingModeProd:
		logger, err = zap.NewProduction()
	default:
		err = fmt.Errorf("invalid logging mode, supported values are 'dev' or 'prod'")
	}
	if err != nil {
		return Config{}, err
	}
	config.Logging.Logger = logger

	return config, nil
}

func (l Logging) IsDev() bool {
	return l.Mode == loggingModeDev
}

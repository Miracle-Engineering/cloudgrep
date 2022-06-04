package cmd

import (
	"errors"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/pflag"
)

type Options struct {
	OutputPath  string
	ConfigPath  string
	Format      bool
	LineNumbers bool
}

func (o *Options) Default() {
	o.Format = true
}

func (o *Options) Bind(flags *pflag.FlagSet) {
	flags.StringVarP(&o.OutputPath, "output-dir", "o", o.OutputPath, "Directory to write generated code to. If not set, writes to stdout")
	flags.StringVarP(&o.ConfigPath, "config", "c", o.ConfigPath, "Location to root config file")
	flags.BoolVar(&o.Format, "format", o.Format, "If true, run generated code through \"go fmt\" before returning")
	flags.BoolVar(&o.LineNumbers, "line-numbers", o.LineNumbers, "If true, append line numbers to each line")
}

func (o *Options) Validate() error {
	var err error

	if !isFile(o.ConfigPath) {
		err = multierror.Append(err, errors.New("--config does not point to a valid config file"))
	}

	return err
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

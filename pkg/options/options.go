package options

import (
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Version  bool   `short:"v" long:"version" description:"Show version"`
	HTTPHost string `long:"bind" description:"HTTP server host"`
	HTTPPort uint   `long:"listen" description:"HTTP server listen port"`
	Prefix   string `long:"prefix" description:"Add a url prefix"`
	SkipOpen bool   `short:"s" long:"skip-open" description:"Skip browser open on start"`
	Config   string `short:"c" long:"config" description:"Configuration file"`
}

var Opts Options

// ParseOptions returns a new options struct from the input arguments
func ParseOptions(args []string) (Options, error) {
	var opts = Options{}

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return opts, err
	}

	if opts.Prefix == "" {
		opts.Prefix = os.Getenv("URL_PREFIX")
	}

	if opts.Prefix != "" && !strings.Contains(opts.Prefix, "/") {
		opts.Prefix = opts.Prefix + "/"
	}
	return opts, nil
}

// SetDefaultOptions parses and assigns the options
func SetDefaultOptions() error {
	opts, err := ParseOptions([]string{})
	if err != nil {
		return err
	}
	Opts = opts
	return nil
}

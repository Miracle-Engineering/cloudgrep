package command

import (
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Debug    bool   `short:"d" long:"debug" description:"Enable debugging mode"`
	HTTPHost string `long:"bind" description:"HTTP server host" default:"localhost"`
	HTTPPort uint   `long:"listen" description:"HTTP server listen port" default:"8080"`
	Prefix   string `long:"prefix" description:"Add a url prefix"`
	SkipOpen bool   `short:"s" long:"skip-open" description:"Skip browser open on start"`
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

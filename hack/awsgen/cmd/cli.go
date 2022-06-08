package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/generator"
	"github.com/run-x/cloudgrep/hack/awsgen/writer"
)

func Run(args []string) {
	err := Do(args)
	if err == nil {
		return
	}

	fmt.Printf("error while running:\n%v", err)
	os.Exit(1)
}

func Do(args []string) error {
	flags := pflag.NewFlagSet("awsgen", pflag.ContinueOnError)

	opts := Options{}
	opts.Default()
	opts.Bind(flags)

	err := flags.Parse(args)
	if errors.Is(err, pflag.ErrHelp) {
		os.Exit(1)
	} else if err != nil {
		return err
	}

	err = opts.Validate()
	if err != nil {
		return err
	}

	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}

	if cfg == nil {
		panic("unexpected nil config")
	}

	err = config.AggregateValidationErrors(cfg.Validate())
	if err != nil {
		return err
	}

	if opts.ValidateOnly {
		return nil
	}

	writer, err := getWriter(opts)
	if err != nil {
		return err
	}

	gen := generator.Generator{
		Format:      opts.Format,
		LineNumbers: opts.LineNumbers,
	}

	err = gen.Generate(writer, *cfg)
	if err != nil {
		return fmt.Errorf("error generating code: %w", err)
	}

	err = writer.Clean()
	if err != nil {
		return fmt.Errorf("cannot clean output: %w", err)
	}

	return nil
}

func getWriter(opts Options) (writer.Writer, error) {
	if opts.OutputPath == "" {
		return writer.NewStreamWriter(os.Stdout), nil
	}

	writer, err := writer.NewDirWriter(opts.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("invalid --output-dir: %w", err)
	}

	return writer, nil
}

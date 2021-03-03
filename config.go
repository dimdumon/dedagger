package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/kaspanet/kaspad/infrastructure/config"
	"github.com/pkg/errors"
)

type configFlags struct {
	dataDir string
	config.NetworkFlags
}

func parseConfig() (*configFlags, error) {
	cfg := &configFlags{}
	parser := flags.NewParser(cfg, flags.HelpFlag)
	parser.Usage = "dedagger [OPTIONS] <DATADIR>"
	remainingArgs, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	err = cfg.ResolveNetwork(parser)
	if err != nil {
		return nil, err
	}

	if len(remainingArgs) == 0 {
		return nil, errors.New("DATADIR is required")
	}

	if len(remainingArgs) > 1 {
		return nil, errors.New("Only 1 positional argument is allowed: DATADIR")
	}
	cfg.dataDir = remainingArgs[0]

	return cfg, nil
}

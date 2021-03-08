package main

import (
	"fmt"
	"os"

	backend "github.com/svarogg/dedagger/backend"
	"github.com/svarogg/dedagger/frontend"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		printErrorAndExit(fmt.Sprintf("Error parsing config: %+v", err))
	}

	be, teardown, err := backend.NewBackend(cfg.dataDir, cfg.ActiveNetParams)
	if err != nil {
		printErrorAndExit(fmt.Sprintf("Error initializing backend: %+v", err))
	}
	defer teardown()

	fe := frontend.NewFrontend(be)
	err = fe.Start()
	if err != nil {
		printErrorAndExit(fmt.Sprintf("Error from frontend: %+v", err))
	}
}

func printErrorAndExit(message string) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", message))
	os.Exit(1)
}

package main

import (
	"fmt"
	"os"

	backend "github.com/svarogg/dedagger/backend"
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

	for _, store := range be.Stores {
		fmt.Println(store.String())
		fmt.Println("=======================")
		for _, method := range store.Methods {
			fmt.Printf("\t%s\n", method)
		}
	}
}

func printErrorAndExit(message string) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", message))
	os.Exit(1)
}

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

	//for _, store := range be.Stores {
	//	fmt.Println(store.String())
	//	fmt.Println("=======================")
	//	for _, method := range store.Methods {
	//		fmt.Printf("\t%s\n", method)
	//	}
	//}

	//hash, err := externalapi.NewDomainHashFromString("fe2b514dbb0cedef5e6d39737bd7f18b38fd197060883fd954dfda4f18e87c42")
	//if err != nil {
	//	printErrorAndExit(fmt.Errorf("error getting hash from string: %+v", err).Error())
	//}
	//outs := be.Call(be.Stores["blockStore"].Methods["Block"], []reflect.Value{reflect.ValueOf(hash)})
	//for i, out := range outs {
	//	json, err := json.Marshal(out.Interface())
	//	if err != nil {
	//		printErrorAndExit(fmt.Errorf("error marshalling output no %d out of %d: %+v", i, len(outs), err).Error())
	//	}
	//	outInterface := out.Interface()
	//	switch outObj := outInterface.(type) {
	//	case *externalapi.DomainHash:
	//		fmt.Printf("%s\n", outObj)
	//	}
	//	fmt.Println(string(json))
	//}
}

func printErrorAndExit(message string) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", message))
	os.Exit(1)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	gate "github.com/handlename/cognito-gate"
)

var (
	version string
)

func main() {
	var (
		versionFlag bool
	)

	flag.BoolVar(&versionFlag, "version", false, "show version")
	flag.Parse()

	if versionFlag {
		fmt.Println("cognito-gate v" + version)
		os.Exit(0)
	}

	if err := gate.Run(os.Getenv("GATE_CONFIG_PATH")); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"

	homed "github.com/gregdel/homed/lib"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var config string
	flag.StringVar(&config, "config", "./config.yaml", "homed configuration file")

	homed, err := homed.New(config)
	if err != nil {
		return err
	}

	return homed.Run()
}

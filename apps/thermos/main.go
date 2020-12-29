package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	debug := false
	var config string
	var mqttBroker string

	flag.StringVar(&config, "config", "~/.homed", "homed configuration directory")
	flag.StringVar(&mqttBroker, "mqttBroker", "tcp://127.0.0.1:1883", "MQTT broker configuration")
	flag.BoolVar(&debug, "debug", false, "enables the debug mode")
	flag.Parse()

	thermos := New(config, mqttBroker, debug)
	if err := thermos.loadRooms(); err != nil {
		return err
	}

	return thermos.Run()
}

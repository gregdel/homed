package main

import (
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
	thermos := New(
		"/home/greg/dev/homed/data/",
		"tcp://192.168.100.50:1883",
	)
	if err := thermos.loadRooms(); err != nil {
		return err
	}

	return thermos.Run()
}

package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	homed "github.com/gregdel/homed/lib"

	_ "github.com/gregdel/homed/lib/components/boiler"
	_ "github.com/gregdel/homed/lib/components/common"
	_ "github.com/gregdel/homed/lib/components/homed_humidity"
	_ "github.com/gregdel/homed/lib/components/homed_temperature"
	_ "github.com/gregdel/homed/lib/components/roller_shutter"
	_ "github.com/gregdel/homed/lib/components/zigbee2mqtt/climate_sensor"
	_ "github.com/gregdel/homed/lib/components/zigbee2mqtt/trv"

	_ "github.com/gregdel/homed/lib/apps/fakehome"
	_ "github.com/gregdel/homed/lib/apps/httpd"
	_ "github.com/gregdel/homed/lib/apps/mqttd"
)

//go:embed build
var embedFS embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var config string
	flag.StringVar(&config, "config", "./config.yaml", "homed configuration file")
	flag.Parse()

	homed, err := homed.New(config, &embedFS)
	if err != nil {
		return err
	}

	return homed.Run()
}

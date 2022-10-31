package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	homed "github.com/gregdel/homed/lib"

	_ "github.com/gregdel/homed/lib/components/boiler"
	_ "github.com/gregdel/homed/lib/components/common"
	_ "github.com/gregdel/homed/lib/components/device_status"
	_ "github.com/gregdel/homed/lib/components/esphome"
	_ "github.com/gregdel/homed/lib/components/homed_temperature"
	_ "github.com/gregdel/homed/lib/components/humidity"
	_ "github.com/gregdel/homed/lib/components/linky"
	_ "github.com/gregdel/homed/lib/components/rtl433"
	_ "github.com/gregdel/homed/lib/components/saswell_trv"
	_ "github.com/gregdel/homed/lib/components/temperature"
	_ "github.com/gregdel/homed/lib/components/tuya_trv"
	_ "github.com/gregdel/homed/lib/components/wifi_signal"
	_ "github.com/gregdel/homed/lib/components/xiaomi_aqara"

	_ "github.com/gregdel/homed/lib/apps/fakehome"
	_ "github.com/gregdel/homed/lib/apps/httpd"
	_ "github.com/gregdel/homed/lib/apps/mqttd"
	_ "github.com/gregdel/homed/lib/apps/rollershutters"
	_ "github.com/gregdel/homed/lib/apps/tempd"
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

	homed, err := homed.New(config, &embedFS)
	if err != nil {
		return err
	}

	return homed.Run()
}

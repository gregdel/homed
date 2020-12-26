package main

import (
	homed "github.com/gregdel/homed/lib"
	"github.com/urfave/cli/v2"
)

var cache *homed.Homed

func init() {
	cache = homed.New("/home/greg/dev/homed/data/")
}

func newApp() *cli.App {
	app := &cli.App{
		Name: "homectl",
		Commands: []*cli.Command{
			handleFileType(homed.FileTypeRoom),
			handleFileType(homed.FileTypeDevice),
			handleFileType(homed.FileTypeSensor),
		},
	}

	return app
}

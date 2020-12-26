package main

import (
	"fmt"

	homed "github.com/gregdel/homed/lib"
	"github.com/kr/pretty"
	"github.com/urfave/cli/v2"
)

func handleFileType(t homed.FileType) *cli.Command {
	return &cli.Command{
		Name:  string(t),
		Usage: "handles " + string(t) + " related stuff",
		Subcommands: []*cli.Command{
			listFiles(t),
			deleteFile(t),
			showFile(t),
			addFile(t),
		},
	}
}

func deleteFile(t homed.FileType) *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "delete a " + string(t),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Name of the " + string(t),
			},
		},
		Action: func(c *cli.Context) error {
			name := c.String("name")
			if name == "" {
				return ErrMissingFileName
			}

			file, err := cache.Load(name, t)
			if err != nil {
				return err
			}

			return cache.Delete(file)
		},
	}
}

func listFiles(t homed.FileType) *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list available " + string(t) + "s",
		Action: func(c *cli.Context) error {
			names, err := cache.ListFiles(t)
			if err != nil {
				return err
			}

			fmt.Println(string(t) + "s:")
			for _, n := range names {
				fmt.Printf("- %s\n", n)
			}

			return nil
		},
	}
}

func showFile(t homed.FileType) *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "show a " + string(t),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Name of the " + string(t),
			},
		},
		Action: func(c *cli.Context) error {
			name := c.String("name")
			if name == "" {
				return ErrMissingFileName
			}

			file, err := cache.Load(name, t)
			if err != nil {
				return err
			}

			pretty.Println(file)
			return nil
		},
	}
}

func addFile(t homed.FileType) *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add a " + string(t),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Name of the " + string(t),
			},
		},
		Action: func(c *cli.Context) error {
			name := c.String("name")
			if name == "" {
				return ErrMissingFileName
			}

			file, err := homed.NewFile(name, t)
			if err != nil {
				return err
			}

			return cache.Save(file)
		},
	}
}

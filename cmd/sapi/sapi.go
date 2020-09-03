package main

import (
	"fmt"
	"os"
	"sapi/cmd/sapi/build"
	"sapi/cmd/sapi/create"
	"sapi/pkg/bootstrap/flag"
)

func main() {
	fs := flag.NewFlagSet()
	fs.Register(
		&flag.BoolFlag{
			Name:   "create",
			Usage:  "--create, Create Service cache file to current dir",
			Action: create.Start,
		},
		&flag.StringFlag{
				Name:   "build",
				Usage:  "--build=NAME, Create service",
				Action: build.Start,
		},
		&flag.StringFlag{
			Name:     "path",
			Usage:    "--path=DIR, Input dir or Output dir",
			Default:  ".",
		},
		)

	err := fs.Parse()
	if err != nil {
		fmt.Println(err)
	} else {
		if len(os.Args) == 1 {
			fs.PrintDefaults()
			os.Exit(0)
		}
	}
}
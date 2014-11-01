package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

//Behold the globals.
var (
	repo Provider
)

const (
	LOG_ERR int = iota
	LOG_WARN
	LOG_INFO
	OFFICIAL_PROVIDER = "s3:https://forklift.microcloud.io"
)

func Log(err error, fatal bool, level int) {
	fmt.Println(err)
	if fatal {
		os.Exit(1)
	}
}
func main() {

	app := cli.NewApp()
	app.Name = "forklift"
	app.Usage = "The practical package manger."

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Be talkative.",
		},
		cli.StringFlag{
			Name:   "provider",
			Value:  OFFICIAL_PROVIDER,
			Usage:  "Repository Address. `type:location`",
			EnvVar: "FORKLIFT_REPO",
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowSubcommandHelp(c)
	}

	app.Commands = []cli.Command{
		build,
		list,
		show,
		clean,
	}

	app.Before = func(c *cli.Context) error {

		var err error
		provider, location := split(c.String("repo"), ":")
		if err != nil {
			Log(err, true, LOG_ERR)
		}

		return nil
	}
	app.Run(os.Args)
}

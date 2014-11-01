package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/forklift/fl/providers"
)

//Behold the globals.
var (
	repo      providers.Provider
	templates = new(template.Template)
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
	app.Usage = "The software provisioning tool."
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
		list,
		versions,
		show,
		build,
		clean,
	}

	app.Before = func(c *cli.Context) error {

		provider, err := providers.NewProvider(c.String("provider"))
		if err != nil {
			Log(err, true, LOG_ERR)
			return err
		}
		repo = *provider
		return nil
	}
	app.Run(os.Args)
}

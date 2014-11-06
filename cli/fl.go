package main

import (
	"os"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/forklift/fl/engine"
	"github.com/forklift/fl/providers"
)

//Behold the globals.
var (
	Log       engine.Logger
	Engine    engine.Engine
	repo      providers.Provider
	templates = new(template.Template)
)

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
			Value:  "s3:https://forklift.microcloud.io",
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

		Log = logrus.New()
		Engine = engine.Engine{Log}

		provider, err := providers.NewProvider(c.String("provider"))
		if err != nil {
			Log.Fatal(err)
		}
		repo = *provider
		return nil
	}
	app.Run(os.Args)
}

package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/forklift/forklift/engine"
	"github.com/forklift/forklift/providers"
)

//Behold the globals.
var (
	Log      engine.Logger
	Engine   *engine.Engine
	Provider providers.Provider
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
		cli.BoolFlag{
			Name:  "robot",
			Usage: "More structure and parsable output.",
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
		Engine = engine.New(Log)

		//TODO: short syntax for default provider! :/location/version.v32.32
		//provider, err := providers.NewProvider(c.String("provider"))
		return nil
	}
	app.Run(os.Args)
}

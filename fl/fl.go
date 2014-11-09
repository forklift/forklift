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
			Usage:  "Default Repo `type:location`",
			EnvVar: "FORKLIFT_REPO",
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowSubcommandHelp(c)
	}

	app.Commands = []cli.Command{
		build,
		clean,
		show,
		versions,
		list,
	}

	app.Before = func(c *cli.Context) error {

		Log = logrus.New()

		err := providers.SetDefault(c.String("provider"))
		if err != nil {
			Log.Error(err)
			return err
		}

		//Fireup a new Engine.
		//INFO: This maybe possible to postpon later
		//      When we actually need it.
		//      Perhaps the package engine could be
		//      just a bunch of functions like go/log.
		Engine = engine.New(Log)
		return nil
	}
	app.Run(os.Args)
}

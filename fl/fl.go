package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/forklift/forklift/engine"
	"github.com/forklift/forklift/providers"
	"github.com/forklift/geppetto/api"
)

//Behold the globals.
var (
	Log      engine.Logger
	Provider providers.Provider
	Server   *api.Server
)

func main() {

	app := cli.NewApp()
	app.Name = "forklift"
	app.Usage = "The software provisioning tool."
	app.Flags = []cli.Flag{

		cli.BoolFlag{
			Name:  "remote",
			Usage: "Run the task using the remote host.",
		},

		cli.StringFlag{
			Name:   "host",
			Usage:  "Remote build server address.",
			Value:  "http://localhost:7070",
			EnvVar: "FORKLIFT_HOST",
		},

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
		ping,
		build,
		clean,
		show,
		versions,
		list,
	}

	app.Before = func(c *cli.Context) error {

		Log = logrus.New()
		engine.Log = Log

		err := providers.SetDefault(c.String("provider"))

		//If no error and remote, connect and ping.
		if err == nil && c.Bool("remote") {
			Server, err = api.NewClient(c.String("host"))
			if err == nil {
				err = Server.Ping()
			}
		}

		if err != nil {
			Log.Error(err)
			return err
		}

		return nil
	}
	app.Run(os.Args)
}

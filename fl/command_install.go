package main

import (
	"github.com/codegangsta/cli"
	"github.com/forklift/forklift/engine"
	"github.com/forklift/forklift/providers"
)

var install = cli.Command{
	Name:   "install",
	Usage:  "Install a package or packages on your system",
	Action: installAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "root",
			Value: "/",
			Usage: "Specify an alternative installation root (default is /).",
		},
	},
}

func installAction(c *cli.Context) {

	arg := c.Args().First()

	if arg == "" {
		cli.ShowSubcommandHelp(c)
		return
	}

	Provider, label, err := providers.Provide(arg)
	if err != nil {
		Log.Fatal(err)
	}

	pack, err := Provider.Fetch(label)
	if err != nil {
		Log.Fatal(err)
	}

	root := c.String("root")
	err = engine.Install(pack, root)
	if err != nil {
		//	Engine.Clean(true)
		//Maybe uninstall here?
	}
}

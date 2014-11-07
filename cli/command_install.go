package main

import (
	"log"

	"github.com/codegangsta/cli"
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
	err := repo.Update()
	if err != nil {
		log.Fatal(err)
	}

	nv, err := repo.Guess(arg)
	if err != nil {
		Log.Fatal(err)
	}

	pack, err := repo.Fetch(nv)
	if err != nil {
		Log.Fatal(err)
	}

	root := c.String("root")
	err = Engine.Install(pack, root)
	if err != nil {
		//	Engine.Clean(true)
		//Maybe uninstall here?
	}
}

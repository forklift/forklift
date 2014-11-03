package main

import "github.com/codegangsta/cli"

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
		Log(err, true, LOG_ERR)
	}

	nv, err := NewNameVersion(arg)
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	pack, err := repo.Fetch(nv.Name, nv.Version)
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	root := c.String("root")
	err := engine.Install(pack, root)
	if err != nil {
		engine.Clean()
	}
}

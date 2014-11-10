package main

import (
	"github.com/codegangsta/cli"
	"github.com/forklift/forklift/engine"
	"github.com/forklift/forklift/providers"
)

var clean = cli.Command{
	Name:   "clean",
	Usage:  "Clean a forklift build env.",
	Action: cleanAction,
}

func cleanAction(c *cli.Context) {

	arg := c.Args().First()

	if arg == "" {
		cli.ShowSubcommandHelp(c)
		return
	}

	provider, label, err := providers.Provide(arg)
	if err != nil {
		Log.Fatal(err)
	}

	location, err := provider.Source(label)
	if err != nil {
		Log.Fatal(err)
	}

	err = engine.Clean(location)
	if err != nil {
		Log.Error(err)
		return
	}
	Log.Info("Clean succesed.")

}

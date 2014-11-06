package main

import (
	"github.com/codegangsta/cli"
	"github.com/forklift/fl/flp"
)

var clean = cli.Command{
	Name:   "clean",
	Usage:  "Clean a forklift build env.",
	Action: cleanAction,
}

func cleanAction(c *cli.Context) {

	pkg, err := flp.ReadPackage()
	if err != nil {
		Log.Fatal(err)
	}

}

package main

import (
	"errors"

	"github.com/codegangsta/cli"
	"github.com/forklift/fl/flp"
)

var clean = cli.Command{
	Name:   "clean",
	Usage:  "Clean a forklift build env.",
	Action: cleanAction,
}

func cleanAction(c *cli.Context) {

	pkg, err := getFileSystemPackage()
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	runClean(pkg)
}

func runClean(pkg *flp.Package) {
	err := runCommands("Cleaning.", pkg.Clean, true)
	if err != nil {
		Log(errors.New("Cleaning faild."), false, LOG_ERR)
	}
}

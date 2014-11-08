package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/forklift/forklift/flp"
	"github.com/forklift/forklift/providers"
)

var build = cli.Command{
	Name:   "build",
	Usage:  "Build a Forklift Package from an florklift.json",
	Action: buildAction,
}

func buildAction(c *cli.Context) {

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

	pkg := path.Join(location, flp.Tag(label.Version))
	//Start creating the package file.
	storage, err := os.Create(pkg)
	if err != nil {
		Log.Error(err)
		return
	}
	defer func() {
		if err != nil {
			os.Remove(pkg)
		}
	}()

	checksum, err := Engine.Build(location, storage)
	if err != nil {
		Log.Error(err)
		return
	}

	Log.Info(fmt.Sprintf("sha256sum: %x ", checksum))
}

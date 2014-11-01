package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

var list = cli.Command{
	Name:   "list",
	Usage:  "Lists all the packages in the index.",
	Action: listAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "files",
			Usage: "Files",
		},
	},
}

func listAction(c *cli.Context) {
	//TODO: Prettify this.

	err := repo.Update()
	if err != nil {
		Log(err, true, 1)
	}

	for name, versions := range repo.Iterator() {
		fmt.Printf(" %-15.15s: %s\n", name, strings.Join(versions, ", "))
	}
}

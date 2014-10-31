package main

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/forklift/fl-go/flp"
)

var clean = cli.Command{
	Name:   "clean",
	Usage:  "Clean a forklift build env.",
	Action: cleanAction,
}

func cleanAction(c *cli.Context) {

	forkliftjson, err := os.Open("forklift.json")
	if err != nil {
		Log(err, true, 1)
	}
	pkg := new(flp.Package)

	err = json.NewDecoder(forkliftjson).Decode(&pkg)
	if err != nil {
		Log(err, true, 1)
	}

	//TODO: Complain about extrenious or missing files.
	//Add support for .forkliftignore

	for _, cmd := range pkg.Clean {
		cmd := exec.Command("sh", "-c", cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil && err != io.EOF {
			Log(err, true, 1)
		}
	}
}

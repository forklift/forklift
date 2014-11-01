package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/forklift/fl/flp"

	"gopkg.in/yaml.v2"
)

var build = cli.Command{
	Name:   "build",
	Usage:  "Build a Forklift Package from an florklift.json",
	Action: buildAction,
}

func buildAction(c *cli.Context) {

	Forkliftfile, err := ioutil.ReadFile("Forkliftfile")
	if err != nil {
		Log(err, true, 1)
	}

	pkg := new(flp.Package)

	err = yaml.Unmarshal(Forkliftfile, &pkg)
	if err != nil {
		Log(err, true, 1)
	}

	//TODO: Complain about extrenious or missing files.
	//Add support for .forkliftignore

	for _, cmd := range pkg.Build {
		cmd := exec.Command("sh", "-c", cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil && err != io.EOF {
			Log(err, true, 1)
		}
	}

	checksum, err := flp.Build(pkg)
	if err != nil {
		Log(err, true, 1)
	}

	fmt.Printf("sha256sum: %s", checksum)
}

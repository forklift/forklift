package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/forklift/fl/flp"
	"gopkg.in/yaml.v2"
)

//TODO: Create a context interface. similar to Package providers
// but for source.
func getFileSystemPackage() (*flp.Package, error) {

	Forkliftfile, err := ioutil.ReadFile("Forkliftfile")
	if err != nil {
		return nil, err
	}

	pkg := new(flp.Package)

	return pkg, yaml.Unmarshal(Forkliftfile, &pkg)

}

func runCommands(step string, cmdlist []string, returnAtFailur bool) error {
	Log(errors.New("Starting: "+step), false, LOG_INFO)

	var err error
	for _, cmd := range cmdlist {
		Log(errors.New("Runnng: "+cmd), false, LOG_INFO)
		cmd := exec.Command("sh", "-c", cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil && err != io.EOF {
			Log(err, false, LOG_ERR)
			if returnAtFailur {
				return err
			}
		}
	}
	return err
}

package main

import "github.com/codegangsta/cli"

var build = cli.Command{
	Name:   "build",
	Usage:  "Build a Forklift Package from an florklift.json",
	Action: buildAction,
}

func buildAction(c *cli.Context) {

	/*
		pkg, err := flp.ReadPackage()
		if err != nil {
			log.Fatal(err)
		}
			err = runCommands("Build.", pkg.Build, true)
			if err != nil {
				Log(errors.New("Buid Faild. Cleaning up."), false, LOG_WARN)
				runClean(pkg)
				log.Fatal(err)
			}

			//TODO: Complain about extrenious or missing files.
			//Add support for .forkliftignore

			checksum, err := flp.Build(pkg)
			if err != nil {
				Log(err, true, LOG_ERR)
			}

			Log(fmt.Errorf("sha256sum: %s", checksum), false, LOG_INFO)
			runClean(pkg)
	*/
}

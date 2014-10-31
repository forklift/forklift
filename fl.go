package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/codegangsta/cli"
	"github.com/forklift/fl-go/flp"
)

type Config struct {
	Cache bool
	R     *url.URL `json:"-"`
}

//Behold the globals.
var (
	config Config
)

const (
	LOG_ERR int = iota
	LOG_WARN
	LOG_INFO
)

func Log(err error, fatal bool, level int) {
	fmt.Println(err)
	if fatal {
		os.Exit(1)
	}
}
func main() {

	app := cli.NewApp()
	app.Name = "forklift"
	app.Usage = "The practical package manger."

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Be talkative.",
		},
	}

	app.Action = func(c *cli.Context) {
		fmt.Println("Main")
	}

	app.Commands = []cli.Command{
		build,
		list,
		show,
		install,
		clean,
	}

	app.Before = func(c *cli.Context) error {

		repo_url := os.Getenv("FORKLIFT_REPO")
		if repo_url == "" {

			Log(fmt.Errorf("Using official repo. %s\n", repo_url), false, LOG_WARN)
			repo_url = flp.OFFICIAL_REGISTRY
		}

		var err error

		config.R, err = url.Parse(repo_url)
		if err != nil {
			Log(err, true, 0)
		}

		return nil
	}
	app.Run(os.Args)
}

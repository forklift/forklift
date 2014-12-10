package main

import "github.com/codegangsta/cli"

var ping = cli.Command{
	Name:  "ping",
	Usage: "Check server. Useful for testing and scripting.",
	Action: func(c *cli.Context) {
		if c.Bool("remote") {
			localPing(c)
		} else {
			remotePing(c)
		}
	},
}

func remotePing(c *cli.Context) {
	err := Server.Ping()
	if err != nil {
		Log.Error(err)
	}

	Log.Print("Pong.")
}

func localPing(c *cli.Context) {

	Log.Error("illegal: Local ping.")
	Log.Print("Did you forgot passing -remote?")
}

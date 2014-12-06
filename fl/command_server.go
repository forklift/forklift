package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/codegangsta/cli"
)

var server = cli.Command{
	Name:   "server",
	Usage:  "Start a Forklift Server deamon.",
	Action: startAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "endpoint",
			Usage: "The address to listen at.",
			Value: "127.0.0.1:5050",
		},
	},
}

var (
	//TODO: Trim the end slashes spaces.
	packagesServerTemplate = ` Listening at : {{.}}
`
)

func startAction(c *cli.Context) {
	//TODO: Prettify this.

	endpoint := c.String("endpoint")

	mux := http.NewServeMux()
	mux.HandleFunc("/_ping", pong)

	log.Printf("Listening at %s", endpoint)

	err := template.Must(template.New("packagesserver").Parse(packagesServerTemplate)).Execute(os.Stdout, endpoint)
	if err != nil {
		Log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(endpoint, mux))
}

func pong(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}

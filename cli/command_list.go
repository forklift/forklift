package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var list = cli.Command{
	Name:   "list",
	Usage:  "Lists all the packages in the index.",
	Action: listAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "filter",
			Usage: "Filter packages.",
		},
	},
}

var (
	//TODO: Trim the end slashes spaces.
	packagesListTemplate = ` Listing Packages From: {{.Location}}{{if .Packages}}
{{range .Packages}}
   {{ . }}{{ end }}{{else}} 
No Packages Found.{{end}}
`
)

func listAction(c *cli.Context) {
	//TODO: Prettify this.

	arg := c.Args().First()

	if arg == "" {
		arg = "*"
	}

	err := Provider.Update()
	if err != nil {
		Log.Fatal(err)
	}

	templates.New("packageslist").Parse(packagesListTemplate)

	err = templates.ExecuteTemplate(os.Stdout, "packageslist", Provider)
	if err != nil {
		Log.Fatal(err)
	}
}

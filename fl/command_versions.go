package main

import (
	"html/template"
	"os"

	"github.com/codegangsta/cli"
)

var versions = cli.Command{
	Name:   "versions",
	Usage:  "Display avaliable versions of a package.",
	Action: versionsAction,
}

var (
	packageVersionsTemplate = ` Listing Package Versions From: {{.Location}}{{if .Packages}}{{range .Versions }}
 {{ . }}{{ end }}{{else}}
 No Packages Found.{{end}}
`
)

func versionsAction(c *cli.Context) {

	arg := c.Args().First()

	if arg == "" {
		cli.ShowSubcommandHelp(c)
		return
	}

	err := Provider.Update()
	if err != nil {
		Log.Fatal(err)
	}

	err = template.Must(template.New("packageversions").Parse(packageVersionsTemplate)).Execute(os.Stdout, Provider)
	if err != nil {
		Log.Fatal(err)
	}
}

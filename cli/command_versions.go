package main

import (
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

	err := repo.Update()
	if err != nil {
		Log.Fatal(err)
	}

	templates.New("packageversions").Parse(packageVersionsTemplate)

	//TODO: Do we need this? does http.FileServer sort? check source.
	// Perhaps we need a semver.Sort interface.
	//sort.Strings(repo.Packages())
	err = templates.ExecuteTemplate(os.Stdout, "packageversions", repo)
	if err != nil {
		Log.Fatal(err)
	}
}

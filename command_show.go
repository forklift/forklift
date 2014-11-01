package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/forklift/flp"
)

var show = cli.Command{
	Name:   "show",
	Usage:  "Display details about a package.",
	Action: showAction,
}

var packageInfoTemplate = `NAME          {{.Name}}
DESCRIPTION   {{.Description}}
VERSION       {{ .Version }} 
KEYWRODS      {{ .Keywrods }} 
HOMEPAGE      {{ .Homepage }} 
BUGS          {{ .Bugs     }} 

AUTHORS       {{ range .Authors }}{{.}}
              {{end}}
OFFICIAL      {{ .Official }} 
PRIVATE       {{ .Private  }} 
LICENSE       {{ .License  }} 

TYPE          {{ .Type }} 
MAIN          {{ .Main }} 
STRUCTURE{{ range .Structure }}
              {{ . }}{{end}}

DEPENDENCIES  {{/* .Dependencies */}} 

INSTALL   
UNINSTALL 
`

func showAction(c *cli.Context) {

	arg := c.Args().First()

	if arg == "" {
		cli.ShowSubcommandHelp(c)
		return
	}
	err := repo.Update()
	if err != nil {
		Log(err, true, 1)
	}

	nv, err := NewNameVersion(arg)
	if err != nil {
		Log(err, true, 1)
	}

	pack, err := repo.Fetch(nv.Name, nv.Version)
	if err != nil {
		Log(err, true, 1)
	}

	pkg, err := flp.Unpack(pack, true)
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	templates.New("packagesinfo").Parse(packageInfoTemplate)

	err = templates.ExecuteTemplate(os.Stdout, "packagesinfo", pkg)
	if err != nil {
		Log(err, true, LOG_ERR)
	}
}

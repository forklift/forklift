package main

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/forklift/fl-go/flp"
	"github.com/hashicorp/go-version"
)

var show = cli.Command{
	Name:   "show",
	Usage:  "Display details about a package or group of packages.",
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

	name := c.Args().First()

	if name == "" {
		cli.ShowSubcommandHelp(c)
		return
	}
	err := GetIndex()
	if err != nil {
		Log(err, true, 1)
	}

	versions, exists := index[name]
	if !exists || len(versions) == 0 {
		Log(fmt.Errorf("Package %s not found.", name), true, 1)
	}

	latest, err := version.Latest(versions)
	if err != nil {
		Log(err, true, 1)
	}

	r := *config.R
	r.Path = path.Join(name, flp.Tag(name, latest))

	pkg, err := flp.Fetch(r, true)
	if err != nil {
		Log(err, true, 1)
	}

	t := template.Must(template.New("packageinfo").Parse(packageInfoTemplate))
	err = t.Execute(os.Stdout, pkg)
	if err != nil {
		Log(err, true, 1)
	}
}

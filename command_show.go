package main

import (
	"fmt"

	"github.com/codegangsta/cli"
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

	fmt.Printf("nv %+v\n", nv)
	latest, err := repo.Get(nv.Name, nv.Version)
	if err != nil {
		Log(err, true, 1)
	}

	fmt.Printf("latest %+v\n", latest)

	//latest, _ := semver.NewVersion("")
	//r := *config.R
	//r.Path = path.Join(arg, flp.Tag(arg, latest))

	/*pkg, err := flp.Fetch(r, true)
	if err != nil {
		Log(err, true, 1)
	}

	t := template.Must(template.New("packageinfo").Parse(packageInfoTemplate))
	err = t.Execute(os.Stdout, pkg)
	if err != nil {
		Log(err, true, 1)
	}*/
}

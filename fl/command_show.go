package main

import (
	"os"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/forklift/forklift/flp"
	"github.com/forklift/forklift/providers"
)

var show = cli.Command{
	Name:   "show",
	Usage:  "Display details about a package.",
	Action: showAction,
}

//Perhaps these can be arranged in more appropriate groups?
var packageInfoTempate = `NAME               {{.Name}}
VERSION            {{ .Version }} 
DESCRIPTION        {{ .Description}}
LICENSE            {{ .License  }} 
KEYWRODS           {{ .Keywrods }} 

PRIVATE            {{ .Private  }} 
REPOSITORY         {{ .Repository }}
BUGS               {{ .Bugs     }} 

OFFICIAL           {{ .Official }} 
MAINTAINERS {{ range .Maintainers }}{{.}}
                   {{end}}

HOMEPAGE           {{ .Homepage }} 
AUTHORS            {{ range .Authors }}{{.}}
                   {{end}}

TYPE               {{ .Type }} 

DEPENDENCIES       {{ range $dep:= .Dependencies}}{{ range $name, $ver := $dep}}{{ $name }} {{ $ver }} 
                   {{ end }} {{ end }}
FILES              {{ range .Files }}{{ . }}
                   {{end}}
INSTALL            {{ range .Install }}{{ . }}
                   {{end}}
UNINSTALL          {{ range .Uninstall }}{{ . }}
                   {{end}}

BUILD DEPENDENCIES {{ range $dep:= .BuildDependencies }}{{ range $name, $ver := $dep}}{{ $name }} {{ $ver }} 
                   {{ end }}{{ end }}

BUILD              {{ range .Build }}{{ . }} 
                   {{end}}

CLEAN              {{ range .Clean }}{{ . }}
                   {{end}}`

func showAction(c *cli.Context) {

	arg := c.Args().First()

	if arg == "" {
		cli.ShowSubcommandHelp(c)
		return
	}

	provider, label, err := providers.Provide(arg)

	if err != nil {
		Log.Fatal(err)
	}

	pack, err := provider.Fetch(label)
	if err != nil {
		Log.Fatal(err)
	}

	pkg, err := flp.Unpack(pack, "", true)
	if err != nil {
		Log.Fatal(err)
	}

	err = template.Must(template.New("packagesinfo").Parse(packageInfoTempate)).Execute(os.Stdout, pkg)
	if err != nil {
		Log.Fatal(err)
	}
}

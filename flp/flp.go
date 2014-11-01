package flp

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/omeid/semver"
)

var (
	ErrInvalidForkliftjson = errors.New("Invalid Forkliftfile")
	ErrMissingForkliftjson = errors.New("No Forkliftfile")
)

func Tag(name string, version *semver.Version) string {
	if path.Ext(name) == ".flp" {
		return name
	}
	return fmt.Sprintf("%s-%s.flp", name, version.String(true))
}

type File struct {
	Meta tar.Header
	Data bytes.Buffer
}

type Package struct {
	Name        string
	Version     string
	Description string

	Keywrods []string

	Private    bool
	Repository string
	Bugs       string

	Official bool
	Authors  []string
	License  string
	Homepage string

	BuildDependencies []string `"yaml:"build-dependencies"`
	Build             []string
	Clean             []string

	Type         string
	Main         string
	Dependencies []string
	Files        []string
	Install      []string
	Uninstall    []string

	//DISCUSS: issues/2
	Runtime struct {
		Kernel       string
		LXC          string
		Libcontainer string
	}

	isStab    bool            `yaml:"-"`
	FilesReal map[string]File `yaml:"-"`
}

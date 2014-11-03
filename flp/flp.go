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
	ErrInvalidForkliftfile = errors.New("Invalid Forkliftfile")
	ErrMissingForkliftfile = errors.New("No Forkliftfile")
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
	License     string
	Description string

	Keywrods []string

	Private    bool
	Repository string
	Bugs       string

	Official    bool
	Maintainers []string

	Homepage string
	Authors  []string

	Type         string
	Main         string
	Dependencies []map[string]string
	Files        []string
	Install      []string
	Uninstall    []string

	BuildDependencies []map[string]string `yaml:"build-dependencies"`
	Build             []string
	Clean             []string

	isStab    bool            `yaml:"-"`
	FilesReal map[string]File `yaml:"-"`
}

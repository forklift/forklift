package flp

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/omeid/semver"
	"gopkg.in/yaml.v2"
)

var (
	ErrInvalidForkliftfile = errors.New("Invalid Forkliftfile")
	ErrMissingForkliftfile = errors.New("No Forkliftfile")
)

func Tag(version *semver.Version) string {
	return fmt.Sprintf("%s-%s.flp", version.Product, version.StringWithMeta())
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
}

func NewPackage(Forkliftfile []byte) (*Package, error) {
	pkg := new(Package)
	return pkg, yaml.Unmarshal(Forkliftfile, &pkg)
}

func ReadPackage() (*Package, error) {
	Forkliftfile, err := ioutil.ReadFile("Forkliftfile")
	if err != nil {
		return nil, err
	}
	return NewPackage(Forkliftfile)
}

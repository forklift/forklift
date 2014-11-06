package providers

import (
	"io"
	"os"
	"path"

	"github.com/forklift/fl/flp"
	"github.com/omeid/semver"
)

func init() {
	List["local"] = &Local{}
}

//				 Name     Versions
type Local struct {
	location string //file system path.
}

func (p *Local) SetLocation(location string) error {
	p.location = location
	return nil
}

func (p Local) Location() string {
	return p.location
}

func (p *Local) Update() error {
	return nil //Error provider Local doesn't support update?
}
func (p *Local) Packages(filter string) ([]string, error) {
	return nil, nil
}

func (p *Local) Versions(product string) ([]string, error) {

	versions := []string{}

	return versions, nil
}

func (p *Local) Fetch(ver *semver.Version) (io.Reader, error) {
	return os.Open(path.Join(p.location, flp.Tag(ver)))
}

func (p *Local) Source(ver *semver.Version) (string, error) {
	return "", nil
}

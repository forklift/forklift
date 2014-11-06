package providers

import (
	"io"
	"os"
	"path"

	"github.com/forklift/fl/flp"
	"github.com/omeid/semver"
)

func init() {
	//List["local"] = &Local{}
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

func (p *Local) Packages() ([]string, error) {
	return nil, nil
}

func (p *Local) Versions() ([]string, error) {

	versions := []string{}

	return versions, nil
}

func (p *Local) Get(name string, ranges string) (*semver.Version, error) {

	versions, err := p.Versions()
	if err != nil {
		return nil, err
	}

	c, err := semver.NewCollection(versions)
	if err != nil {
		return nil, err
	}

	return c.Latest(ranges)
}

func (p *Local) Fetch(ver *semver.Version) (io.Reader, error) {
	return os.Open(path.Join(p.location, flp.Tag(ver)))
}

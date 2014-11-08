package providers

import (
	"io"
	"os"
	"path"

	"github.com/forklift/forklift/flp"
	"github.com/omeid/semver"
)

func init() {
	List["local"] = &Local{}
}

//				 Name     Versions
type Local struct{}

func (p *Local) Parse(label string) (*Label, error) {

	pkg, err := flp.ReadPackage(label)
	if err != nil {
		return nil, err
	}

	ver, err := semver.NewVersion(pkg.Name + "-" + pkg.Version)
	if err != nil {
		return nil, err
	}

	return &Label{label, ver, false}, nil
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

func (p *Local) Fetch(l *Label) (io.Reader, error) {
	return os.Open(path.Join(l.Location, flp.Tag(l.Version)))
}

func (p *Local) Source(l *Label) (string, error) {
	return l.Location, nil
}

package providers

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/forklift/forklift/flp"
	"github.com/forklift/forklift/semver"
)

func init() {
	list["local"] = &local{}
}

//				 Name     Versions
type local struct{}

func (p *local) Parse(labelstring string) (*Label, error) {

	Label := &Label{}
	var err error

	if path.Ext(labelstring) == ".flp" {
		Label.Location = path.Dir(labelstring)
		ver := strings.TrimRight(path.Base(labelstring), path.Ext(labelstring))
		Label.Version, err = semver.NewVersion(ver)

		if err != nil {
			return Label, err
		}
	} else {

		Label.Location = labelstring

		pkg, err := flp.ReadPackage(labelstring)
		if err != nil {
			return nil, err
		}

		Label.Version, err = semver.NewVersion(pkg.Name + "-" + pkg.Version)
		if err != nil {
			return nil, err
		}
	}

	return Label, nil
}

func (p *local) SetLocation(string) error {
	return nil //Local provider does not support location.
}

func (p *local) Update() error {
	return nil //Error provider Local doesn't support update?
}
func (p *local) Packages(filter string) ([]string, error) {
	return nil, errors.New("Provider `Local` doesn't support Package listing.")
}

func (p *local) Versions(product string) ([]string, error) {

	versions := []string{}

	return versions, nil
}

func (p *local) Fetch(l *Label) (io.Reader, error) {
	return os.Open(path.Join(l.Location, flp.Tag(l.Version)))
}

func (p *local) Source(l *Label) (string, error) {
	return l.Location, nil
}

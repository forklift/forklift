package providers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/forklift/fl/flp"
	"github.com/omeid/semver"
)

func init() {
	List["go"] = &GO{}
}

//				 Name     Versions
type GO struct {
	location *url.URL
	packages []string
}

func (p *GO) SetLocation(location string) error {
	var err error
	p.location, err = url.Parse(location)
	return err
}

func (p GO) Location() string {
	return p.location.String()
}

func (p *GO) Update() error {

	if p.location == nil {
		return errors.New("Provider unset.")
	}
	return nil
}

func (p *GO) Packages(filter string) ([]string, error) {

	//If no filter or catch all, return it all.
	if filter == "*" || filter == "" {
		return p.packages, nil
	}

	//If there is no globing or filter charchters,
	//There can be only one matching package.
	one := strings.IndexAny(filter, "*?[") < 0

	filtered := []string{}

	for _, pkg := range p.packages {

		matched, err := filepath.Match(filter, pkg)
		if err != nil {
			return filtered, err
		}
		if matched {
			filtered = append(filtered, pkg)

			if one {
				break
			}
		}
	}

	return filtered, nil
}

func (p *GO) Versions(filter string) ([]string, error) {

	versions := []string{}

	if p.location == nil {
		return versions, errors.New("Location is not set.")
	}

	u := *p.location

	u.Path = path.Join(u.Path, filter) + "/"

	//versions, err :=  p.Packages(filter)
	versions = nil
	var err error

	if err != nil {
		return versions, err
	}

	return versions, nil
}

func (p *GO) Get(name string, ranges string) (*semver.Version, error) {

	versions, err := p.Versions(name)
	if err != nil {
		return nil, err
	}

	c, err := semver.NewCollection(versions)
	if err != nil {
		return nil, err
	}

	return c.Latest(ranges)
}

func (p *GO) Fetch(ver *semver.Version) (io.Reader, error) {

	if p.location == nil {
		return nil, errors.New("Package not found.")
	}

	u := *p.location
	u.Path = path.Join(u.Path, ver.Product, flp.Tag(ver))

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		return nil, fmt.Errorf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))
	}
	return res.Body, nil
}

func (p *GO) Source(ver *semver.Version) (string, error) {
	return "", nil
}

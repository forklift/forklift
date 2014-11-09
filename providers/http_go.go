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

	"github.com/forklift/forklift/flp"
)

func init() {
	list["go"] = &http_go{}
}

//				 Name     Versions
type http_go struct {
	location *url.URL
	packages []string
}

func (p *http_go) SetLocation(location string) error {
	var err error
	//TODO: Should we Ping already?
	p.location, err = url.Parse(location)
	return err
}

func (p *http_go) Update() error {
	if p.location == nil {
		return errors.New("Provider unset.")
	}
	return nil
}

func (p *http_go) Parse(label string) (*Label, error) {
	return nil, nil
}

func (p *http_go) Packages(filter string) ([]string, error) {

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

func (p *http_go) Versions(filter string) ([]string, error) {

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

func (p *http_go) Fetch(l *Label) (io.Reader, error) {

	if l.Location == "" {
		return nil, errors.New("Package not found.")
	}

	if l.Version == nil || l.Version.Product == "" {
		return nil, nil // ErrorMissingProduct
	}

	location := path.Join(l.Location, l.Version.Product, flp.Tag(l.Version))

	res, err := http.Get(location)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		return nil, fmt.Errorf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))
	}
	return res.Body, nil
}

func (p *http_go) Source(l *Label) (string, error) {
	return "", nil
}

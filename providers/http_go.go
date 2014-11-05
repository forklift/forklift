package providers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/forklift/fl/flp"
	"github.com/omeid/semver"
)

func init() {
	//	List["go"] = &GO{}
}

//				 Name     Versions
type GO struct {
	location *url.URL
	f        string //The filter. TODO: Remove it.

	index struct {
		XMLNAME  xml.Name `xml:"pre"`
		Packages []string `xml:"a"`
	}
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

	var err error
	r := "^[a-zA-Z0-9].*/$"

	if p.f != "" && p.f != "*" {
		r = p.f + ".*/"
	}

	reg, err := regexp.Compile(r)
	if err != nil {
		return err
	}

	f := func(v []byte, o *string) bool {
		s := string(v)
		*o = strings.TrimRight(s, "/")
		return reg.MatchString(s)
	}
	p.index.Packages, err = getXML(p.location.String(), "a", f)

	if err != nil {
		return err
	}

	return nil
}

func (p *GO) SetFilter(f string) {
	p.f = f
}

func (p *GO) Packages() []string {
	return p.index.Packages
}

func (p *GO) Versions(filter string) ([]string, error) {

	versions := []string{}

	if p.location == nil {
		return versions, errors.New("Location is not set.")
	}

	u := *p.location

	u.Path = path.Join(u.Path, filter) + "/"

	reg, err := regexp.Compile("^" + filter + "(-[0-9].*)?$")
	if err != nil {
		return versions, err
	}

	f := func(v []byte, o *string) bool {
		s := string(v)
		s = strings.TrimRight(s, "\\.flp")
		_, err := semver.Parts(s)
		if err != nil {
			return false
		}
		*o = s
		return reg.MatchString(s)
	}
	versions, err = getXML(u.String(), "a", f)

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

func (p *GO) Fetch(ver semver.Version) (io.Reader, error) {

	if p.location == nil {
		return nil, errors.New("Package not found.")
	}

	latest, err := p.Get(name, ranges)
	if err != nil {
		return nil, err
	}

	u := *p.location
	u.Path = path.Join(u.Path, name, flp.Tag(name, latest))

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

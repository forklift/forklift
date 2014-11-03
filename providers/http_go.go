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
	"github.com/forklift/fl/util"
	"github.com/omeid/semver"
)

func init() {
	List["go"] = &GO{}
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

	f := func(v string, o *string) bool {
		*o = strings.TrimRight(v, "/")
		return reg.MatchString(v)
	}
	p.index.Packages, err = util.GetHTMLElements(p.location.String(), "a", "href", f)

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

	f := func(v string, o *string) bool {
		v = strings.TrimRight(v, "\\.flp")
		_, err := semver.Parts(v)
		if err != nil {
			return false
		}

		//*o = strings.TrimLeft(v, p.f+"-")
		*o = v //For easier copypasta and scripting.
		return reg.MatchString(v)
	}
	versions, err = util.GetHTMLElements(u.String(), "a", "href", f)

	if err != nil {
		return versions, err
	}

	return versions, nil
}

func (p *GO) Get(name string, ranges string) (*semver.Version, error) {

	p.SetFilter(name)
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

func (p *GO) Fetch(name string, ranges string) (io.Reader, error) {

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

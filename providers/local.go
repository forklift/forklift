package providers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/forklift/fl/flp"
	"github.com/forklift/fl/util"
	"github.com/omeid/semver"
)

func init() {
	//List["local"] = &Local{}
}

//				 Name     Versions
type Local struct {
	f string //The filter. TODO: Remove it.

	index struct {
		XMLNAME  xml.Name `xml:"pre"`
		Packages []string `xml:"a"`
	}
}

func (p *Local) SetLocation(location string) error {
	return nil
}

func (p Local) Location() string {
	return "Not Location Support."
}

func (p *Local) Packages() []string {
	return p.index.Packages
}

func (p *Local) Versions() ([]string, error) {

	versions := []string{}

	if p.location == nil {
		return versions, errors.New("Package not found.")
	}

	u := *p.location

	u.Path = path.Join(u.Path, p.f) + "/"

	reg, err := regexp.Compile("^" + p.f + "(-[0-9].*)?$")
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

func (p *Local) Get(name string, ranges string) (*semver.Version, error) {

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

func (p *Local) Fetch(name string, ranges string) (io.Reader, error) {

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

package providers

import (
	"encoding/xml"
	"net/url"
)

func init() {
	List["http"] = new(HTTP)
}

//				 Name     Versions
type HTTP struct {
	provider
	location *url.URL

	Index struct {
		XMLNAME  xml.Name `xml:"pre"`
		packages []string
	}
}

func (p HTTP) SetLocation(location string) error {
	var err error
	p.location, err = url.Parse(location)
	return err
}

func (p HTTP) Update() error {

	/*
			repo := path.Join(config.R.String(), "Forkliftindex")

			res, err := http.Get(repo)
			if err != nil {
				return err
			}

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))
			}

			err = json.NewDecoder(res.Body).Decode(&index)
			if err != nil {
				return err
			}
		}
	*/
	return nil
}

func (p HTTP) Get(filter string) []string {
	return []string{}
}

/*
func Fetch(path url.URL, stabOnly bool) (*Package, error) {
	res, err := http.Get(path.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))
	}

	return Unpack(res.Body, stabOnly)
}
*/

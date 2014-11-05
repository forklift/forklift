package providers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

//Doesn't work with void (self closing) elements.
// The `element` body is passed to `filter` and will be added to return slice
// only if it returns true.

func getXML(url string, element string, filter func([]byte, *string) bool) ([]string, error) {

	s := []string{}

	resp, err := http.Get(url)
	if err != nil {
		return s, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s, fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	decoder := xml.NewDecoder(resp.Body)

	for token, err := decoder.Token(); token != nil; token, err = decoder.Token() {
		if err != nil {
			if err != io.EOF {
				return err
			}
			return s, nil
		}

		switch e := token.(type) {
		case xml.StartElement:
			if e.Name.Local == element {
				token, err := decoder.Token() //Expecting xml.CharData
				if err != nil {
					return err, nil
				}
				if e, ok := token.(xml.CharData); ok {
					s := ""
					if filter(xml.CharData, &s) {
						s = append(s, v)
					}
				} // we don't handle the else case for a more graceful parser.
			}
		}
	}
	return s, nil
}

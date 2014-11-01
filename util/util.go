package util

import (
	"fmt"
	"net/http"

	"code.google.com/p/go.net/html"
)

func GetHTMLElements(url string, element string, key string, filter func(string, *string) bool) ([]string, error) {
	// result

	s := []string{}

	resp, err := http.Get(url)
	if err != nil {
		return s, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s, fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	//TODO: Consider XPATH instead of walking manually.

	//Build the filter.
	// parser
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return s, err
	}

	// The walker.
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		var v string
		if n.Type == html.ElementNode && n.Data == element {
			for _, a := range n.Attr {
				if a.Key == key && filter(a.Val, &v) {
					s = append(s, v)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	//Lets walk the doc.
	walk(doc)

	// done
	return s, nil
}

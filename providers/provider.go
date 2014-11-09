package providers

import (
	"errors"
	"io"
	"strings"

	"github.com/forklift/forklift/semver"
)

var list = map[string]Provider{}

var (
	ErrorProviderNoSuch  = errors.New("[Provider] No Such Provider.")
	ErrorProviderInvalid = errors.New("[Provider] Provider Invalid.")
	ErrorProviderMissing = errors.New("[Provider] Please specifiy a provider.")

	ErrorLabelEmpty   = errors.New("[Provider] The Label is empty.")
	ErrorLabelInvalid = errors.New("[Provider] Invalid Label. Require 'provider:location package'")
)

type Label struct {
	Location string
	Version  *semver.Version

	//If it is not the exact version requested for.
	Alt bool
}

func SetDefault(uri string) error {

	if uri == "" {
		return ErrorProviderInvalid
	}

	parts := strings.SplitN(uri, ":", 2)

	//Do we have such provider?
	provider, ok := list[parts[0]]
	if !ok {
		return ErrorProviderNoSuch
	}

	// If there is no location,
	// lets hope it accepts an empty location.
	if len(parts) < 2 {
		parts[1] = ""
	}

	list["default"] = provider

	return list["default"].SetLocation(parts[1])
}

func Provide(uri string) (Provider, *Label, error) {
	if uri == "" {
		return nil, nil, ErrorLabelEmpty
	}

	var (
		provider    Provider
		labelstring string
	)

	first := uri[0]

	//File system?
	if first == '.' || first == '/' {
		provider = &local{}
		labelstring = uri

	}

	//Remote provider?
	if provider == nil {
		parts := strings.SplitN(uri, ":", 2)

		if len(parts) == 2 {

			var ok bool
			if provider, ok = list[parts[0]]; !ok {
				return nil, nil, ErrorProviderNoSuch
			}
			labelstring = parts[1]

		}
	}

	//Lets try if there is a default one.
	if provider == nil {
		var ok bool
		if provider, ok = list["default"]; !ok {
			return nil, nil, ErrorProviderMissing
		}
		labelstring = uri
	}

	label, err := provider.Parse(labelstring)
	return provider, label, err

}

type Provider interface {

	//This method sets the location when it applies.
	//Return an error on invalid location, ignore if it
	//doesn't applies.
	SetLocation(string) error
	//This method parses a Label
	//into Locatoin and Version, returns an error if
	//invalid Label.
	Parse(string) (*Label, error)

	//Used for remote providers to update their repository list.
	//Providers MUST attempt to reconnect and only return an error if connection
	//fails after reasonable number of retrys.
	Update() error

	//List all packages, accepts a filter.
	Packages(string) ([]string, error)

	//List all version of a package, accepts a product name.
	Versions(string) ([]string, error)

	//Fetches a specific package.
	Fetch(*Label) (io.Reader, error)

	// Provides the location for the source of a specific version.
	// The location must be a location on the local file system, preferably under /tmp/ when possible.
	// It is the stevenhong: responsiblity of the provider to fetch and extract the source.
	Source(*Label) (string, error)
}

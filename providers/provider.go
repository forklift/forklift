package providers

import (
	"errors"
	"io"
	"strings"

	"github.com/omeid/semver"
)

var List = map[string]Provider{}

var (
	ErrorPackageProviderMissing = errors.New("Package Provider Missing.")
	ErrorPackageProviderInvlaid = errors.New("Package Provider Invalid.")
	ErrorNoSuchProvider         = errors.New("No Such Provider.")
)

func NewProvider(pp string) (*Provider, error) {
	if pp == "" {
		return nil, ErrorPackageProviderMissing
	}

	p := strings.SplitN(pp, ":", 2)
	if len(p) < 2 {
		return nil, ErrorPackageProviderInvlaid
	}

	if provider, ok := List[p[0]]; ok {
		err := provider.SetLocation(p[1])
		return &provider, err
	}

	return nil, ErrorNoSuchProvider
}

type Provider interface {

	//If applies, this funciton should set the address of provider
	//Return error on invalid format only, if provider doesn't requires
	//an address or it can't be set, return nil.
	SetLocation(string) error

	//Return the location of provier, if it doesn't applies, return an
	//Appropriate message. (i.e: "Local" Provider does not support location.")
	Location() string

	//Used for remote providers to update their repository list.
	//Providers MUST attempt to reconnect and only return an error if connection
	//fails after reasonable number of retrys.
	Update() error

	//List all packages, accepts a filter.
	Packages(string) ([]string, error)

	//List all version of a package.
	Versions(string) ([]string, error)

	//Fetches a specific package.
	Fetch(*semver.Version) (io.Reader, error)

	// Provides the location for the source of a specific version.
	// The location must be a location on the local file system, preferably under /tmp/ when possible.
	// It is the responsiblity of the provider to fetch and extract the source.
	Source(*semver.Version) (string, error)
}

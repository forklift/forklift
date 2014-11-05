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

	//Used for remote providers to check network
	//Providers MUST attempt to reconnect and only return an error if connection
	//fails after reasonable number of retrys.
	Ping() error

	//List all packages, accets a filter.
	Packages(string) []string

	//List all version of a package.
	Versions(string) ([]string, error)

	//Fetches a specific package.
	Fetch(*semver.Version) (io.Reader, error)

	// Provides the location for the source of a specific version.
	// The location must be a location on the local file system, preferably under /tmp/ when applies.
	// Thus, the provider may need to fetch, and extract if neccessary.
	Source(*semver.Version) (string, error)
}

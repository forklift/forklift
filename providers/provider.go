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
	SetLocation(string) error
	Location() string

	Ping() bool                               //Used for remote providers to check network
	Packages(string) []string                 //List all pacages, accepts a filter.
	Versions(string) ([]string, error)        //List all versions of a package.
	Fetch(*semver.Version) (io.Reader, error) //Fetches a specific package.

	// Provides the location for the source of a specific version
	// Fetch, and extract if neccessary.
	Source(*semver.Version) (string, error)
}

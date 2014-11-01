package providers

import (
	"errors"
	"strings"

	"github.com/forklift/fl/flp"
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

	Update() error

	//TODO:There should be away to get ride of the SetFilter method
	// And pass the filter directly to Packagees and Version methods.
	SetFilter(string)
	Packages() []string
	Versions() ([]string, error)
	Get(string, string) (*semver.Version, error) //Accept package name and a range provide the best option, empty if no matching version.

	flp.Fetcher
}

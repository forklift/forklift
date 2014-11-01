package providers

import (
	"errors"
	"strings"
)

var List = map[string]Provider{}

var (
	ErrorPackageProviderMissing = errors.New("Package Provider Missing.")
	ErrorPackageProviderInvlaid = errors.New("Package Provider Invalid.")
	ErrorNoSuchProvider         = errors.New("No Such Provider.")
)

func GetProvider(pp string) (*Provider, error) {
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
	Iterator() map[string][]string

	SetLocation(string) error
	Update() error
	Get(string) []string
}

type provider struct {
	Index map[string][]string
}

func (p *provider) Iterator() map[string][]string { return p.Index }

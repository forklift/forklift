package semver

import "sort"

var (
	ErrorVersionNotFound = Error{"Version %s not found.", ""}
)

// Collection is a type that implements the sort.Interface interface
// so that versions can be sorted.
type Collection []*Version

func NewCollection(raws []string) (Collection, error) {

	versions := make(Collection, len(raws))

	for i, raw := range raws {
		v, err := NewVersion(raw)
		if err != nil {
			return versions, err
		}
		versions[i] = v
	}

	return versions, nil
}

func (vc Collection) Len() int {
	return len(vc)
}

//TODO: Check why sort.Reverse isn't working.
func (vc Collection) Less(i, j int) bool {
	return vc[j].LessThan(vc[i])
}

func (vc Collection) Swap(i, j int) {
	vc[i], vc[j] = vc[j], vc[i]
}

func (vc Collection) Latest(constraints string) (*Version, error) {

	//This puts the highest versions on top.
	sort.Sort(vc)
	c, err := NewConstraint(constraints)
	if err != nil {
		return nil, err
	}

	for _, v := range vc {
		if c.Check(v) {
			return v, nil
		}
	}

	return nil, ErrorVersionNotFound.Fault(constraints)
}

package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	//Exported vars.
	ErrorMalformedVersionPre = Error{"Malformed Version Pre part: %s", ""}
	PrereleaseTags           = []string{"pre", "alpha", "beta", "rc"}
)

var (
	preRegexpRaw = `(` + strings.Join(PrereleaseTags, "|") + `)?(?:.([0-9]+))?`
	preRegexp    = regexp.MustCompile(preRegexpRaw)
)

type Pre struct {
	Ident string
	Patch int
}

func NewPre(raw string) (*Pre, error) {

	pre := new(Pre)

	parts := preRegexp.FindStringSubmatch(raw)
	if parts == nil {
		return nil, ErrorMalformedVersionPre
	}

	pre.Ident = parts[1]

	if len(parts) > 1 && parts[2] != "" {
		//We pass the raw through regex, so
		// parts[2] is either empty, or a digit.
		pre.Patch, _ = strconv.Atoi(parts[2])
	}

	return pre, nil
}

func (p *Pre) String() string {
	return fmt.Sprintf("%s.%d", p.Ident, p.Patch)
}

var stages = map[string]int{
	"pre":   0,
	"alpha": 1,
	"beta":  2,
	"rc":    3,
}

// Compares Pre v to o:
// -1 == v is less than o
// 0 == v is equal to o
// 1 == v is greater than o

func (v *Pre) Compare(o *Pre) int {

	if stages[v.Ident] > stages[o.Ident] {
		return 1
	}

	if stages[v.Ident] < stages[o.Ident] {
		return -1
	}

	if v.Patch > o.Patch {
		return 1
	}

	if v.Patch < o.Patch {
		return -1
	}

	return 0
}

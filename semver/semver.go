package semver

import (
	"fmt"
	"regexp"
)

var (
	ErrorMalformedVersion = Error{"Malformed Version: %s", ""}

	versionRegexpRaw = `(?:([a-zA-Z0-9]+[a-zA-Z0-9-]*[a-zA-Z0-9]+)-)?` +
		`([0-9]+(?:\.[0-9]+){0,2})` +
		`(?:\-(` + preRegexpRaw + `))?` + //The extra capture group is to pass the whole release.tag to NewPre.
		`(?:\+([0-9A-Za-z.]+))?`
	versionRegexp = regexp.MustCompile(`^` + versionRegexpRaw + `$`)
)

type Version struct {
	Raw string

	Product string

	Major uint16
	Minor uint16
	Patch uint16
	Pre   *Pre
	Meta  string
}

func Parts(raw string) ([]string, error) {
	parts := versionRegexp.FindStringSubmatch(raw)
	if parts == nil {
		return nil, ErrorMalformedVersion.Fault(raw)
	}
	return parts, nil
}

func NewVersion(raw string) (*Version, error) {

	parts, err := Parts(raw)
	if err != nil {
		return nil, err
	}

	ver := new(Version)

	ver.Raw = parts[0]
	ver.Product = parts[1]

	fmt.Sscanf(parts[2], "%d.%d.%d", &ver.Major, &ver.Minor, &ver.Patch)

	pre, err := NewPre(parts[3])
	if err != nil {
		return nil, err
	}

	if parts[3] != "" {
		ver.Pre = pre
	}
	ver.Meta = parts[6]
	return ver, nil
}

func (v *Version) String() string {

	ver := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	pre := v.Pre.String()
	if pre != "" {
		ver = ver + "-" + pre
	}

	return ver
}

func (v *Version) StringWithMeta() string {
	ver := v.String()

	if v.Meta != "" {
		ver = ver + "+" + v.Meta
	}
	return ver

}

// Compares Versions v to o:
// -1 == v is older than o
// 0 == v is same as o
// 1 == v is newer than o

func (v *Version) Compare(o *Version) int {

	//The MIT License

	if v.Major > o.Major {
		return 1
	}

	if v.Major < o.Major {
		return -1
	}

	if v.Minor > o.Minor {
		return 1
	}

	if v.Minor < o.Minor {
		return -1
	}

	if v.Patch > o.Patch {
		return 1
	}

	if v.Patch < o.Patch {
		return -1
	}

	if v.Pre == nil {
		return 1
	}

	if o.Pre == nil {
		return -1
	}

	//Both has prereleases, compare them.
	return v.Pre.Compare(o.Pre)
}

func (v *Version) LessThan(o *Version) bool {
	return v.Compare(o) == -1
}
func (v *Version) GreaterThan(o *Version) bool {
	return v.Compare(o) == 1
}
func (v *Version) Equal(o *Version) bool {
	return v.Compare(o) == 0
}
func (v *Version) NotEqual(o *Version) bool {
	return v.Compare(o) != 0
}
func (v *Version) GreaterThanOrEqual(o *Version) bool {
	return v.Compare(o) != -1
}
func (v *Version) LessThanOrEqual(o *Version) bool {
	return v.Compare(o) != 1
}

package semver

import (
	"testing"
)

type testset []struct {
	pass bool // This flag indicts if this set of vers pass as valid or not.
	vers []string
}

func TestNewVersion(t *testing.T) {

	tests := testset{
		//This is a set of valid versions.
		{
			true,
			[]string{"package-ver-1.2.3", "1.2.3", "1.2", "1", "1.2-beta.5", "1.2.0-beta", "1.2.0-alpha+metadata", "1.2.0-beta.1", "1.2.0-alpha+metadata.3"},
		},
		{
			false,
			[]string{"package-ver-1.2.3+foo-dat", "1-0", "1.2.3-unicorn", "foo", "1.2-5", "1.2.0-x.Y.0+metadata", "1.2.0-0alpha+metadata", "1.2.0-0", "1.2.0+sd+fa"},
		},
	}

	for _, test := range tests {

		for _, ver := range test.vers {
			v, err := NewVersion(ver)

			if (err == nil) != test.pass {
				t.Fatalf("Test Failed.\n '%s'.\n Expected Error: %t\n ERROR: %v\n VERSION: %v\v", ver, !test.pass, err, v)
			}
		}

	}
}

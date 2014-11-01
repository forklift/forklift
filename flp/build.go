package flp

import (
	"encoding/hex"
	"errors"
	"os"
	"path"

	"github.com/omeid/semver"
)

// Build. Building a version package from Forkliftfile.
func Build(pkg *Package) (checksum string, err error) {

	ver, err := semver.NewVersion(pkg.Version)
	if err != nil {
		return pkg.Version, errors.New("Invalid Package name or version.")
	}

	//TODO: Complain about extrenious or missing files.
	//Add support for .forkliftignore

	tag := Tag(pkg.Name, ver)
	//Start creating the package file.
	flpfile, err := os.Create(tag)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			os.Remove(tag)
		}
	}()

	//hash := sha256.New()
	//w := io.MultiWriter(flpfile, hash)
	// Add the forklift.json to the start.
	for i, p := range pkg.Files {
		pkg.Files[i] = path.Join("root", p)
	}

	pkg.Files = append([]string{"Forkliftfile"}, pkg.Files...)
	//Pack it.
	sum, err := Pack(pkg.Files, flpfile)
	checksum = hex.EncodeToString(sum)
	if err != nil {
		return "", err
	}

	return checksum, nil
}

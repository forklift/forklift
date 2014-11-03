package flp

import "os"

// Build. Building a version package from Forkliftfile.
func Build(pkg *Package) (checksum string, err error) {

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
}

package flp

import (
	"archive/tar"
	"errors"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Unpack(pack io.Reader, stabOnly bool) (*Package, error) {

	tar := tar.NewReader(pack)

	hdr, err := tar.Next()

	//TODO: Bettter error handling.
	if err != nil || err == io.EOF || hdr.Name != "Forkliftfile" {
		return nil, ErrMissingForkliftfile
	}

	pkg := new(Package)
	Forkliftfile, err := ioutil.ReadAll(tar)
	if err != nil {
		return nil, errors.New("Error reading Forklift file from tar. This should neverh happen. Please open an issue at github.com/forklift/fl/issues")
	}

	err = yaml.Unmarshal(Forkliftfile, &pkg)
	if err != nil {
		return nil, ErrInvalidForkliftfile
	}

	if stabOnly {
		return pkg, nil
	}

	pkg.FilesReal = make(map[string]File)
	for {
		hdr, err := tar.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return nil, err
		}

		f := new(File)
		f.Meta = *hdr

		_, err = io.Copy(&f.Data, tar)

		if err != nil {
			return nil, err
		}
		pkg.FilesReal[hdr.Name] = *f
	}
	return pkg, nil
}

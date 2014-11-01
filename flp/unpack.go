package flp

import (
	"archive/tar"
	"encoding/json"
	"io"
)

func Unpack(pack io.Reader, stabOnly bool) (*Package, error) {

	tar := tar.NewReader(pack)

	hdr, err := tar.Next()

	//TODO: Bettter error handling.
	if err != nil || err == io.EOF || hdr.Name != "forklift.json" {
		return nil, ErrMissingForkliftjson
	}

	pkg := new(Package)
	if err := json.NewDecoder(tar).Decode(&pkg); err != nil {
		return nil, ErrInvalidForkliftjson
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

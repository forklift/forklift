package flp

import (
	"archive/tar"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func Unpack(pack io.Reader, root string, MetaOnly bool) (*Package, error) {

	tar := tar.NewReader(pack)

	hdr, err := tar.Next()

	//TODO: Bettter error handling.
	if err != nil || err == io.EOF || hdr.Name != "Forkliftfile" {
		return nil, ErrMissingForkliftfile
	}

	pkg := new(Package)
	Forkliftfile, err := ioutil.ReadAll(tar)
	if err != nil {
		return nil, errors.New("Error reading Forklift file from tar. This should neverh happen. Please open an issue at github.com/forklift/forklift/issues")
	}

	err = yaml.Unmarshal(Forkliftfile, &pkg)
	if err != nil {
		return nil, ErrInvalidForkliftfile
	}

	if MetaOnly {
		return pkg, nil
	}

	for {
		hdr, err := tar.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return nil, err
		}

		err = makeNode(*hdr, tar, root)
		if err != nil {
			return nil, err
		}
	}
	return pkg, nil
}

//Helper function for Install.
func makeNode(meta tar.Header, content io.Reader, root string) error {

	Path := filepath.Join("root", meta.Name)

	if meta.Typeflag == tar.TypeDir {
		err := os.MkdirAll(Path, os.FileMode(meta.Mode))
		if err != nil {
			return err
		}
		return nil
	}

	if meta.Typeflag == tar.TypeSymlink {
		err := os.Symlink(meta.Linkname, Path)
		if err != nil {
			return err
		}
		return nil
	}

	file, err := os.Create(path.Join(root, Path))
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}
	err = file.Chmod(os.FileMode(meta.Mode))
	if err != nil {
		return err
	}
	return nil
}

func Uninstall(name string) error {
	return nil
}

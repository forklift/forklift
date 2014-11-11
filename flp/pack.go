package flp

import (
	"archive/tar"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/forklift/forklift/semver"
)

func Pack(root string, pkg *Package, storage io.WriteCloser) (checksum []byte, err error) {

	hash := sha256.New()
	w := io.MultiWriter(storage, hash)

	tar := tar.NewWriter(w)

	//Taking care of closing everything.
	defer func() {
		if e := tar.Close(); e != nil && err == nil {
			err = e
		}
		if e := storage.Close(); e != nil && err == nil {
			err = e
		}
		checksum = hash.Sum(nil)
	}()

	_, err = semver.NewVersion(pkg.Version)
	if err != nil {
		return nil, errors.New("Invalid Package name or version.")
	}

	for i, p := range pkg.Files {
		pkg.Files[i] = path.Join("root", p)
	}

	//Add Forkliftfile as first File without the root prefix.
	pkg.Files = append([]string{"Forkliftfile"}, pkg.Files...)

	//Add files to the package.
	for _, path := range pkg.Files {
		err = writeFile(root, path, tar)
		if err != nil {
			return nil, err
		}
	}

	checksum = hash.Sum(nil)
	if err != nil {
		return nil, err
	}
	return checksum, nil
}

func writeFile(root string, name string, tarstream *tar.Writer) error {

	longpath := filepath.Join(root, name)
	file, err := os.Open(longpath)
	if err != nil {
		return err
	}
	defer file.Close()

	fi, err := os.Lstat(longpath)
	if err != nil {
		return err
	}

	var link string

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, err = filepath.EvalSymlinks(longpath)
		if err != nil {
			return err
		}
	}

	th, err := tar.FileInfoHeader(fi, link)
	if err != nil {
		return fmt.Errorf("Header Error: %q: %v", name, err)
	}

	th.Name = name

	if err := tarstream.WriteHeader(th); err != nil {
		return fmt.Errorf("Write Failed: %q: %v", name, err)
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return nil
	}

	if !fi.IsDir() {
		if _, err := io.Copy(tarstream, file); err != nil {
			return fmt.Errorf("Write Failed: %q: %v", name, err)
		}
		return nil
	}

	subs, err := file.Readdirnames(0)
	if err != nil {
		return fmt.Errorf("Read Failed: %q: %v", name, err)
	}
	for _, s := range subs {
		if err := writeFile(root, filepath.Join(name, s), tarstream); err != nil {
			return err
		}
	}
	return nil
}

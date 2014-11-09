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

func Pack(pkg *Package, storage io.WriteCloser) (checksum []byte, err error) {

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
		err = writeFile(path, tar)
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

func writeFile(p string, tarstream *tar.Writer) error {

	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	fi, err := os.Lstat(p)
	if err != nil {
		return err
	}

	var link string

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, err = filepath.EvalSymlinks(p)
		if err != nil {
			return err
		}
	}

	th, err := tar.FileInfoHeader(fi, link)
	if err != nil {
		return fmt.Errorf("Header Error: %q: %v", p, err)
	}

	th.Name = p

	if err := tarstream.WriteHeader(th); err != nil {
		return fmt.Errorf("Write Failed: %q: %v", p, err)
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return nil
	}

	if !fi.IsDir() {
		if _, err := io.Copy(tarstream, file); err != nil {
			return fmt.Errorf("Write Failed: %q: %v", p, err)
		}
		return nil
	}

	subs, err := file.Readdirnames(0)
	if err != nil {
		return fmt.Errorf("Read Failed: %q: %v", p, err)
	}
	for _, name := range subs {
		if err := writeFile(filepath.Join(p, name), tarstream); err != nil {
			return err
		}
	}
	return nil
}

package flp

import (
	"archive/tar"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Pack(files []string, storage io.WriteCloser) (checksum []byte, err error) {

	hash := sha256.New()
	w := io.MultiWriter(storage, hash)

	tar := tar.NewWriter(w)

	defer func() {
		if e := tar.Close(); e != nil && err == nil {
			err = e
		}
		if e := storage.Close(); e != nil && err == nil {
			err = e
		}
		checksum = hash.Sum(nil)
	}()

	for _, path := range files {
		err = writeFile(path, tar)
		if err != nil {
			return nil, err
		}
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

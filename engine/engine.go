package engine

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/forklift/fl/flp"
)

var Log Logger

type Engine struct {
}

//Build
func (e Engine) Build(dir string, storage io.WriteCloser) ([]byte, error) {

	bounce, err := bouncer(dir)
	defer bounce()
	if err != nil {
		return nil, err
	}

	pkg, err := flp.ReadPackage()
	if err != nil {
		return nil, err
	}

	err = runCommands("Build.", pkg.Build, true)
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return flp.Pack(pkg, storage)
}

//Clean

func (e Engine) Clean(dir string, storage io.WriteCloser) error {

	bounce, err := bouncer(dir)
	defer bounce()
	if err != nil {
		return nil, err
	}

	pkg, err := flp.ReadPackage()
	if err != nil {
		return nil, err
	}

	//It runCommands with false never returns anything,
	// All the errors are logged directly.
	return runCommands("Cleaning.", pkg.Clean, false)
}

// Install
func (e Engine) Install(pack io.Reader, root string) error {

	if root != "/" {
		//Not using the default root, so we need to chroot (Change root).
		//This requires root user or sudo access.
		//Can probably be fix with fakeroot/fakechroot.
		err = syscall.Chroot(root)
	} else {
		//All post comand installs are run from / of filesystem.
		bounce, err := bouncer("/")
		defer bounce()
		if err != nil {
			return err
		}
	}

	if err != nil {
		Log.Error(err)
		return err
	}

	pkg, err := flp.Unpack(pack, false)
	if err != nil {
		Log.Error(err)
		return err
	}

	for _, file := range pkg.FilesReal {
		err := makeNode(file.Meta, &file.Data, root)
		if err != nil {
			Log(err, true, LOG_ERR) //Clean up here.
		}
	}

	err = runCommands("Post Install", pkg.Install, true)
	if err != nil {
		Log(errors.New("Post install Faild. Uninstalling."), false, LOG_WARN)
		runClean(pkg)
		Log(err, true, LOG_ERR) //We can die now. :/
	}

	Log.Print("Package installed successfuly.", pkg.Version)
	return nil
}

//Helper function for Install.
func makeNode(meta tar.Header, content io.Reader, root string) error {

	Path := filepath.Join(root, meta.Name)

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

	file, err := os.Create(Path)
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
}

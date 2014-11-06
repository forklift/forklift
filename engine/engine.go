package engine

import (
	"io"
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

	err = run("Build.", pkg.Build, true)
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
		return err
	}

	pkg, err := flp.ReadPackage()
	if err != nil {
		return err
	}

	//It runCommands with false never returns anything,
	// All the errors are logged directly.
	return run("Cleaning.", pkg.Clean, false)
}

// Install
func (e Engine) Install(pack io.Reader, root string) error {

	var err error
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

	pkg, err := flp.Unpack(pack, root, false)
	if err != nil {
		Log.Error(err)
		return err
	}

	err = run("Post Install", pkg.Install, true)
	if err != nil {
		Log.Error(err)
		//Log.Warn("Post install Faild. Uninstalling.")
		//e.Uninstall(pkg)
	}

	Log.Print("Package installed successfuly.", pkg.Version)
	return nil
}

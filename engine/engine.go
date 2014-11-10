package engine

import (
	"io"
	"syscall"

	"github.com/forklift/forklift/flp"
)

var Log Logger

//Build
func Build(dir string, storage io.WriteCloser) ([]byte, error) {

	pkg, err := flp.ReadPackage(dir)
	if err != nil {
		return nil, err
	}

	bounce, err := bouncer(dir)
	defer bounce()
	if err != nil {
		return nil, err
	}

	Log.Info("Starting: build...")
	err = run(Log, pkg.Build, true)
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return flp.Pack(pkg, storage)
}

//Clean

func Clean(dir string) error {

	pkg, err := flp.ReadPackage(dir)
	if err != nil {
		return err
	}

	bounce, err := bouncer(dir)
	defer bounce()
	if err != nil {
		return err
	}

	//"run" with false never returns anything,
	// All the errors are logged directly.
	Log.Info("Starting: Cleaning..")
	return run(Log, pkg.Clean, false)
}

// Install
func Install(pack io.Reader, root string) error {

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

	Log.Info("Starting: Post Install..")
	err = run(Log, pkg.Install, true)
	if err != nil {
		Log.Error(err)
		//e.log.Warn("Post install Faild. Uninstalling.")
		//e.Uninstall(pkg)
	}

	Log.Print("Package installed successfuly.", pkg.Version)
	return nil
}

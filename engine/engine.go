package engine

import (
	"io"
	"syscall"

	"github.com/forklift/fl/flp"
)

func New(log Logger) *Engine {
	return &Engine{log: log}
}

type Engine struct {
	log Logger
}

//Build
func (e *Engine) Build(dir string, storage io.WriteCloser) ([]byte, error) {

	bounce, err := bouncer(dir)
	defer bounce()
	if err != nil {
		return nil, err
	}

	pkg, err := flp.ReadPackage()
	if err != nil {
		return nil, err
	}

	err = run(e.log, "Build.", pkg.Build, true)
	if err != nil {
		e.log.Error(err)
		return nil, err
	}
	return flp.Pack(pkg, storage)
}

//Clean

func (e *Engine) Clean(dir string, storage io.WriteCloser) error {

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
	return run(e.log, "Cleaning.", pkg.Clean, false)
}

// Install
func (e *Engine) Install(pack io.Reader, root string) error {

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
		e.log.Error(err)
		return err
	}

	pkg, err := flp.Unpack(pack, root, false)
	if err != nil {
		e.log.Error(err)
		return err
	}

	err = run(e.log, "Post Install", pkg.Install, true)
	if err != nil {
		e.log.Error(err)
		//e.log.Warn("Post install Faild. Uninstalling.")
		//e.Uninstall(pkg)
	}

	e.log.Print("Package installed successfuly.", pkg.Version)
	return nil
}

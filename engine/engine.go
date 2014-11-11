package engine

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/forklift/forklift/flp"
)

var Log Logger

//A helper function that runs a slice of cmd.
func run(root string, dir string, cmdlist []string, returnAtFailur bool, log Logger) error {

	for _, cmd := range cmdlist {
		log.Info("Running: ", cmd)
		cmd := exec.Command("sh", "-c", cmd)
		if root != "" {
			cmd.SysProcAttr = &syscall.SysProcAttr{Chroot: root}
		}
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil && err != io.EOF {
			log.Warn(err)
			if returnAtFailur {
				return err
			}
		}
	}
	return nil
}

//Build
func Build(dir string, storage io.WriteCloser) ([]byte, error) {

	pkg, err := flp.ReadPackage(dir)
	if err != nil {
		return nil, err
	}

	Log.Info("Starting: build...")
	err = run("", dir, pkg.Build, true, Log)
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return flp.Pack(dir, pkg, storage)
}

//Clean

func Clean(dir string) error {

	pkg, err := flp.ReadPackage(dir)
	if err != nil {
		return err
	}

	//"run" with false never returns anything,
	// All the errors are logged directly.
	Log.Info("Starting: Cleaning..")
	return run("", dir, pkg.Clean, false, Log)
}

// Install
func Install(pack io.Reader, root string) error {

	pkg, err := flp.Unpack(pack, root, false)
	if err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Starting: Post Install..")
	err = run(root, "/", pkg.Install, true, Log)
	if err != nil {
		Log.Error(err)
		//e.log.Warn("Post install Faild. Uninstalling.")
		//e.Uninstall(pkg)
	}

	Log.Print("Package installed successfuly.", pkg.Version)
	return nil
}

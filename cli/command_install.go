package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/forklift/fl/flp"
)

var install = cli.Command{
	Name:   "install",
	Usage:  "Install a package or packages on your system",
	Action: installAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "root",
			Value: "/",
			Usage: "Specify an alternative installation root (default is /).",
		},
	},
}

func installAction(c *cli.Context) {

	arg := c.Args().First()

	if arg == "" {
		cli.ShowSubcommandHelp(c)
		return
	}
	err := repo.Update()
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	nv, err := NewNameVersion(arg)
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	pack, err := repo.Fetch(nv.Name, nv.Version)
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	//err = engine.Install(pack, root)

	pkg, err := flp.Unpack(pack, false)
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	root := c.String("root")
	for _, file := range pkg.FilesReal {
		err := makeNode(file.Meta, &file.Data, root)
		if err != nil {
			Log(err, true, LOG_ERR) //Clean up here.
		}
	}

	if root != "/" {
		//Not using the default root, so we need to chroot (Change root).
		//This requires root user or sudo access.
		//Can probably be fix with fakeroot/fakechroot.
		err = syscall.Chroot(root)
	} else {
		//All post comand installs are run from / of filesystem.
		err = os.Chdir("/")
	}
	if err != nil {
		Log(err, true, LOG_ERR)
	}

	err = runCommands("Post Install", pkg.Install, true)
	if err != nil {
		Log(errors.New("Post install Faild. Uninstalling."), false, LOG_WARN)
		runClean(pkg)
		Log(err, true, LOG_ERR) //We can die now. :/
	}

	Log(fmt.Errorf("Package %s installed successfuly.", pkg.Version), false, 2)
}

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

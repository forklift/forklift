package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
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
	/*
		name := c.Args().First()

		if name == "" {
			cli.ShowSubcommandHelp(c)
			return
		}
		versions := make([]*version.Version, len(versionsRaw))
		for i, raw := range versionsRaw {
			v, _ := version.NewVersion(raw)
			versions[i] = v
		}

		sort.Sort(version.Collection(versions))

		latest := versions[0]

		r := *config.R
		r.Path = path.Join(name, flp.Tag(name, latest))

		pkg, err := flp.Fetch(r, false)
		if err != nil {
			Log(err, true, LOG_ERR)
		}
		for _, file := range pkg.FilesReal {
			err := makeNode(file.Meta, &file.Data, c.String("root"))
			if err != nil {
				Log(err, true, LOG_ERR) //Clean up here.
			}
		}

		Log(fmt.Errorf("Package %s installed successfuly.", name), false, 2)
	*/
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

package engine

import (
	"io"
	"os"
	"os/exec"
)

//Bouncer changes the directory and returns a "returning" function which
//when called, it will return you to the previous directory.
func bouncer(dir string) (func() error, error) {

	//The infamous nope function, it does nothing.
	//We return the nope on failur so it can be called
	// on Defer regardless of the error status.
	nope := func() error { return nil }

	pwd, err := os.Getwd()
	if err != nil {
		//Couldn't even get the current working directory. No bounc required.
		return nope, err
	}

	//We have the current working dirctory.
	//build the bounce function.
	bounc := func() error { return os.Chdir(pwd) }

	err = os.Chdir(dir)
	if err != nil {
		//Couldn't Change Directory, so return a nope, no bounc required.
		return nope, err
	}

	return bounc, err
}

func run(log Logger, step string, cmdlist []string, returnAtFailur bool) error {

	log.Info("Starting: ", step)

	for _, cmd := range cmdlist {
		log.Info("Starting: ", cmd)
		cmd := exec.Command("sh", "-c", cmd)
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

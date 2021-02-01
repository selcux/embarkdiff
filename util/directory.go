package util

import (
	"os"
	"os/user"
)

//HomeDirectory gives the home directory of the current user
func HomeDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

//DirExists returns true if the given path is a directory
func DirExists(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return false, err
	}

	return !os.IsNotExist(err) && info.IsDir(), nil
}

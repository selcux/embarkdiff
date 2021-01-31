package util

import (
	"os"
	"os/user"
)

func HomeDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func DirExists(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return false, err
	}

	return !os.IsNotExist(err) && info.IsDir(), nil
}

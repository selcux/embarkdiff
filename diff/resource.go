package diff

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const resourceFile string = "embarkdiff.json"

type Resource struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func (r *Resource) Validate() bool {
	return r.Source != "" && r.Target != ""
}

func (r *Resource) Write() error {
	buffer, err := json.Marshal(r)
	if err != nil {
		return err
	}

	resFile, err := resourcePath(resourceFile)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(resFile, buffer, 0644)
}

func Read() (*Resource, error) {
	resFile, err := resourcePath(resourceFile)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(resFile)
	if os.IsNotExist(err) {
		return new(Resource), nil
	}

	buffer, err := ioutil.ReadFile(resFile)
	if err != nil {
		return nil, err
	}

	resource := new(Resource)
	err = json.Unmarshal(buffer, resource)
	return resource, err
}

func homeDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func resourcePath(resFile string) (string, error) {
	homeDir, err := homeDirectory()
	if err != nil {
		return "", err
	}

	return path.Join(homeDir, resourceFile), nil
}

func dirExists(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return false, err
	}

	return os.IsExist(err) && info.IsDir(), nil
}

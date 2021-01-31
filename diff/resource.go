package diff

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const resourceFile string = "embarkdiff.json"

type ResourceInterface interface {
	Validate() bool
	Write() error
	Load() error

	Source() string
	SetSource(dir string) error

	Target() string
	SetTarget(dir string) error
}

type Resource struct {
	source string
	target string
}

func NewResource() ResourceInterface {
	return &Resource{}
}

func (r *Resource) Source() string {
	return r.source
}

func (r *Resource) SetSource(dir string) error {
	return setDir(dir, &r.source)
}

func (r *Resource) Target() string {
	return r.target
}

func (r *Resource) SetTarget(dir string) error {
	return setDir(dir, &r.target)
}

func (r *Resource) Validate() bool {
	return r.source != "" && r.target != ""
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

func (r *Resource) Load() error {
	resFile, err := resourcePath(resourceFile)
	if err != nil {
		return err
	}

	_, err = os.Stat(resFile)
	if os.IsNotExist(err) {
		return nil
	}

	buffer, err := ioutil.ReadFile(resFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buffer, r)
	return err
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

func setDir(dir string, resField *string) error {
	if resField == nil {
		return errors.New("resource field cannot be nil")
	}

	exists, err := dirExists(dir)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("directory %s does not exists", dir)
	}

	*resField = dir

	return nil
}

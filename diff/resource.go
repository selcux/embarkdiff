package diff

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/selcux/embarkdiff/util"
	"io/ioutil"
	"os"
	"path"
)

//resourceFile refers to the name of the state file
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

//Resource is the structure to keep state of command line inputs
type Resource struct {
	source string
	target string
}

func NewResource() ResourceInterface {
	return &Resource{}
}

//Source retrieves the source directory path
func (r *Resource) Source() string {
	return r.source
}

//SetSource updates the source directory path
func (r *Resource) SetSource(dir string) error {
	return setDir(dir, &r.source)
}

//Target retrieves the target directory path
func (r *Resource) Target() string {
	return r.target
}

//SetTarget updates the target directory path
func (r *Resource) SetTarget(dir string) error {
	return setDir(dir, &r.target)
}

//Validate the directory paths for not being empty
func (r *Resource) Validate() bool {
	return r.source != "" && r.target != ""
}

//Write to the resource file
func (r *Resource) Write() error {
	buffer, err := json.Marshal(r)
	if err != nil {
		return err
	}

	resFile, err := r.resourcePath()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(resFile, buffer, 0644)
}

//Load the resource file into `Resource` struct
func (r *Resource) Load() error {
	resFile, err := r.resourcePath()
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

//resourcePath retrieves the path of the state file
func (r *Resource) resourcePath() (string, error) {
	homeDir, err := util.HomeDirectory()
	if err != nil {
		return "", err
	}

	return path.Join(homeDir, resourceFile), nil
}

//setDir updates the given field with the given directory path
func setDir(dir string, resField *string) error {
	if resField == nil {
		return errors.New("resource field cannot be nil")
	}

	exists, err := util.DirExists(dir)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("directory %s does not exists", dir)
	}

	*resField = dir

	return nil
}

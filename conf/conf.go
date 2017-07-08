package conf

import (
	"errors"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/skibish/ddns/do"
)

// Configuration is a structure of the configuration
type Configuration struct {
	Token   string                 `yaml:"token"`
	Domain  string                 `yaml:"domain"`
	Records []do.Record            `yaml:"records"`
	Notify  map[string]interface{} `yaml:"notify"`
	Params  map[string]string      `yaml:"params"`
}

// valid checks that provided configuration is valid
func (c *Configuration) valid() error {
	if c.Token == "" {
		return errors.New("Token can't be empty")
	}

	if c.Domain == "" {
		return errors.New("Domain can't be empty")
	}

	return nil
}

// NewConfiguration read configuration file
// and return *Configuration
func NewConfiguration(path string) (*Configuration, error) {
	path = os.ExpandEnv(path)

	file, errRead := ioutil.ReadFile(path)
	if errRead != nil {
		return nil, errRead
	}
	var cf Configuration
	errUn := yaml.Unmarshal(file, &cf)
	if errUn != nil {
		return nil, errUn
	}

	errValid := cf.valid()
	if errValid != nil {
		return nil, errValid
	}

	if cf.Params == nil {
		cf.Params = map[string]string{}
	}

	return &cf, nil
}

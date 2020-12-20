package conf

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func createTmpFile(t *testing.T) (string, func()) {
	t.Helper()

	f, err := ioutil.TempFile("", "demo-*.yml")
	if err != nil {
		t.Errorf("Failed to create temp file")
	}

	rm := func() {
		os.Remove(f.Name())
	}
	return f.Name(), rm
}

func TestNewConfigurationMultipleDomainsSuccess(t *testing.T) {
	fname, rm := createTmpFile(t)
	defer rm()

	errWrite := ioutil.WriteFile(fname, []byte(`token: amazing
domains:
  - example.com
  - example.net
records:
  - type: A
    name: www`), 0644)

	if errWrite != nil {
		t.Error("Failed to write file")
		return
	}

	conf, errConf := NewConfiguration(fname)
	if errConf != nil {
		t.Errorf("Got error: %s", errConf.Error())
		return
	}

	if len(conf.Domains) != 2 {
		t.Errorf("Expected two domains in the list, got %v", len(conf.Domains))
		return
	}

	if conf.Domains[0] != "example.com" {
		t.Errorf("Expected example.com, got %s", conf.Domains[0])
		return
	}

	if conf.Domains[1] != "example.net" {
		t.Errorf("Expected example.net, got %s", conf.Domains[0])
		return
	}

	if conf.Records[0].Name != "www" {
		t.Errorf("Expected www, got %s", conf.Records[0].Name)
		return
	}
}

func TestNewConfigurationReadFail(t *testing.T) {
	filePath := "/tmp/demo1.yml"

	_, err := NewConfiguration(filePath)
	if err == nil {
		t.Errorf("Everything is OK, but should be error: %v", err)
		return
	}
}

func TestNewConfigurationParseError(t *testing.T) {
	fname, rm := createTmpFile(t)
	defer rm()

	errWrite := ioutil.WriteFile(fname, []byte(`is not yml`), 0644)

	if errWrite != nil {
		t.Error("Failed to write file")
		return
	}

	_, errConf := NewConfiguration(fname)
	if !strings.Contains(errConf.Error(), "yaml: unmarshal errors") {
		t.Error("Should be error, but everything is OK")
		return
	}
}

func TestNewConfigurationValid(t *testing.T) {
	fname, rm := createTmpFile(t)
	defer rm()

	// check for token
	errWrite := ioutil.WriteFile(fname, []byte(`token: ""
domains:
  - example.com`), 0644)

	if errWrite != nil {
		t.Error("Failed to write file")
		return
	}
	_, errConf := NewConfiguration(fname)
	if errConf.Error() != "token can't be empty" {
		t.Error("Should be error, but everything is OK")
		return
	}

	// check for domains
	errWrite3 := ioutil.WriteFile(fname, []byte(`token: abc
domains: [""]`), 0644)

	if errWrite3 != nil {
		t.Error("Failed to write file")
		return
	}
	_, errConf3 := NewConfiguration(fname)
	if errConf3.Error() != "domains can't be empty" {
		t.Error("Should be error, but everything is OK")
		return
	}
}

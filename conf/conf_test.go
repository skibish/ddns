package conf

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewConfigurationSuccess(t *testing.T) {
	filePath := "/tmp/demo.yml"
	defer os.Remove(filePath)

	errWrite := ioutil.WriteFile(filePath, []byte(`token: amazing
domain: example.com
records:
  - type: A
    name: www`), 0644)

	if errWrite != nil {
		t.Error("Failed to write file")
		return
	}

	conf, errConf := NewConfiguration(filePath)
	if errConf != nil {
		t.Errorf("Got error: %s", errConf.Error())
		return
	}

	if conf.Domain != "example.com" {
		t.Errorf("Expected example.com, got %s", conf.Domain)
		return
	}

	if conf.Records[0].Name != "www" {
		t.Errorf("Expected www, got %s", conf.Records[0].Name)
		return
	}
}

func TestNewConfigurationReadFail(t *testing.T) {
	filePath := "/tmp/demo1.yml"

	_, errConf := NewConfiguration(filePath)
	if errConf.Error() != "open /tmp/demo1.yml: no such file or directory" {
		t.Error("Got error, but should be OK")
		return
	}
}

func TestNewConfigurationParseError(t *testing.T) {
	filePath := "/tmp/demo.yml"
	defer os.Remove(filePath)

	errWrite := ioutil.WriteFile(filePath, []byte(`is not yml`), 0644)

	if errWrite != nil {
		t.Error("Failed to write file")
		return
	}

	_, errConf := NewConfiguration(filePath)
	if !strings.Contains(errConf.Error(), "yaml: unmarshal errors") {
		t.Error("Should be error, but everything is OK")
		return
	}
}

func TestNewConfigurationValid(t *testing.T) {
	filePath := "/tmp/demo.yml"
	defer os.Remove(filePath)

	// check for token
	errWrite := ioutil.WriteFile(filePath, []byte(`token: ""
domain: example.com`), 0644)

	if errWrite != nil {
		t.Error("Failed to write file")
		return
	}
	_, errConf := NewConfiguration(filePath)
	if errConf.Error() != "token can't be empty" {
		t.Error("Should be error, but everything is OK")
		return
	}

	// check for domain
	errWrite2 := ioutil.WriteFile(filePath, []byte(`token: abc
domain: ""`), 0644)

	if errWrite2 != nil {
		t.Error("Failed to write file")
		return
	}
	_, errConf2 := NewConfiguration(filePath)
	if errConf2.Error() != "domain can't be empty" {
		t.Error("Should be error, but everything is OK")
		return
	}
}

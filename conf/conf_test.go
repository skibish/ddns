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

func TestNewConfigurationMultipleDomainsSuccess(t *testing.T) {
	filePath := "/tmp/demo.yml"
	defer os.Remove(filePath)

	errWrite := ioutil.WriteFile(filePath, []byte(`token: amazing
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

	conf, errConf := NewConfiguration(filePath)
	if errConf != nil {
		t.Errorf("Got error: %s", errConf.Error())
		return
	}

	if len(conf.Domains) != 2 {
		t.Errorf("Expected two domains in the list, got %v", len(conf.Domains))
		return
	}

	if conf.Domains[0] != "example.com" {
		t.Errorf("Expected example.com, got %s", conf.Domain)
		return
	}

	if conf.Domains[1] != "example.net" {
		t.Errorf("Expected example.net, got %s", conf.Domain)
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

	// check for domains
	errWrite3 := ioutil.WriteFile(filePath, []byte(`token: abc
domains: [""]`), 0644)

	if errWrite3 != nil {
		t.Error("Failed to write file")
		return
	}
	_, errConf3 := NewConfiguration(filePath)
	if errConf3.Error() != "domain can't be empty" {
		t.Error("Should be error, but everything is OK")
		return
	}
}

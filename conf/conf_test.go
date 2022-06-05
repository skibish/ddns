package conf

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func createTmpFile(t *testing.T) (string, func()) {
	is := is.New(t)
	is.Helper()

	f, err := ioutil.TempFile("", "demo-*.yml")
	is.NoErr(err)

	rm := func() {
		os.Remove(f.Name())
	}
	return f.Name(), rm
}

func TestNewConfigurationMultipleDomainsSuccess(t *testing.T) {
	is := is.New(t)
	fname, rm := createTmpFile(t)
	defer rm()

	err := ioutil.WriteFile(fname, []byte(`token: amazing
domains:
  example.com:
    - type: A
      name: www
  example.net:
    - type: A
      name: www`), 0644)

	is.NoErr(err)

	conf, err := NewConfiguration(fname)
	is.NoErr(err)

	is.Equal(len(conf.Domains), 2)

	_, ok := conf.Domains["example.com"]
	is.True(ok)

	_, ok = conf.Domains["example.net"]
	is.True(ok)
	is.Equal(conf.Domains["example.com"][0].Name, "www")
}

func TestNewConfigurationReadFail(t *testing.T) {
	is := is.New(t)
	_, err := NewConfiguration("/tmp/demo1.yml")
	if err == nil {
		is.Fail() // should be error because no such file
	}
}

func TestNewConfigurationParseError(t *testing.T) {
	is := is.New(t)
	fname, rm := createTmpFile(t)
	defer rm()

	err := ioutil.WriteFile(fname, []byte(`is not yml`), 0644)
	is.NoErr(err)

	_, err = NewConfiguration(fname)
	is.True(strings.Contains(err.Error(), "yaml: unmarshal errors"))
}

func TestNewConfigurationValid(t *testing.T) {
	is := is.New(t)
	fname, rm := createTmpFile(t)
	defer rm()

	// check for token
	err := ioutil.WriteFile(fname, []byte(`token: ""
domains:
  example.com:`), 0644)
	is.NoErr(err)

	_, err = NewConfiguration(fname)
	is.True(strings.Contains(err.Error(), "token can't be empty"))

	// check for domains
	err = ioutil.WriteFile(fname, []byte(`token: abc
domains:`), 0644)
	is.NoErr(err)

	_, err = NewConfiguration(fname)
	is.True(strings.Contains(err.Error(), "domains can't be empty"))

	// check for domain records
	err = ioutil.WriteFile(fname, []byte(`token: abc
domains:
  example.com: []`), 0644)
	is.NoErr(err)

	_, err = NewConfiguration(fname)
	is.True(strings.Contains(err.Error(), "records can't be empty"))
}

func TestEnvVarsAreRead(t *testing.T) {

	is := is.New(t)
	fname, rm := createTmpFile(t)
	defer rm()

	err := ioutil.WriteFile(fname, []byte(`domains:
  example.com:
    - type: A
      name: www`), 0644)

	is.NoErr(err)

	os.Setenv("DDNS_TOKEN", "abc")
	os.Setenv("DDNS_CHECKPERIOD", "60s")
	os.Setenv("DDNS_REQUESTTIMEOUT", "12s")
	os.Setenv("DDNS_IPV6", "true")
	conf, err := NewConfiguration(fname)
	is.NoErr(err)

	is.Equal("abc", conf.Token)
	is.Equal(60*time.Second, conf.CheckPeriod)
	is.Equal(12*time.Second, conf.RequestTimeout)
	is.Equal(true, conf.IPv6)
}

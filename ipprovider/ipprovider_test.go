package ipprovider

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

type fakeProviderOne struct{}

func (f fakeProviderOne) GetIP() (string, error) {
	return "", errors.New("No IP found")
}

type fakeProviderTwo struct{}

func (f fakeProviderTwo) GetIP() (string, error) {
	return "45.45.45.45", nil
}

func TestGetIP(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	log.SetOutput(bw)

	i := New()

	i.Register(&fakeProviderOne{}, &fakeProviderTwo{})

	ip := i.GetIP()
	if ip != "45.45.45.45" {
		t.Error("aaa")
		return
	}

	errFlush := bw.Flush()
	if errFlush != nil {
		t.Errorf("Error flushing log message: %q", errFlush.Error())
		return
	}

	if !strings.Contains(b.String(), "level=warning") {
		t.Error("There should be level=warning in log")
		return
	}

}

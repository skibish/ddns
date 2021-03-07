package notifier

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGotifyInit(t *testing.T) {
	_, errConv := initGotifyNotifier("cfg")
	if errConv == nil {
		t.Error("Should be conversion error, but got nothing")
		return
	}

	m := make(map[interface{}]interface{})
	m["app_url"] = 123
	_, errValue := initGotifyNotifier(m)
	if errValue == nil {
		t.Errorf("Should be error, because value is not a string")
		return
	}

	m = make(map[interface{}]interface{})
	m[123] = "app_url"
	_, errKey := initGotifyNotifier(m)
	if errKey == nil {
		t.Errorf("Should be error, because key is not a string")
		return
	}

	m = make(map[interface{}]interface{})
	m["app_url"] = "https://gotify.example.com/"
	g, errUnexpected := initGotifyNotifier(m)
	if errUnexpected != nil {
		t.Error(errUnexpected)
		return
	}

	if g.AppURL != "https://gotify.example.com" {
		t.Errorf("app_url should be %q, but got %q", "https://gotify.example.com", g.AppURL)
		return
	}

}

func TestGotifyFire(t *testing.T) {

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	g := GotifyConfig{
		AppURL:   "https://example.com/",
		AppToken: "1234",
		errorW:   writer,
	}

	g.send = func(msg string) error {
		return nil
	}

	entry := logrus.Entry{Message: "awesome message", Level: logrus.InfoLevel}
	errFire := g.Fire(&entry)
	if errFire != nil {
		t.Errorf("handled an error on Gotify request: %q", errFire.Error())
		return
	}

	// test, got some error from Gotify
	g.send = func(string) error {
		return errors.New("something went wrong")
	}

	errFire = g.Fire(&entry)
	if errFire == nil {
		t.Error("Should be error, but it's OK")
		return
	}

	// test, DEBUG is ignored
	entry = logrus.Entry{Message: "ignored", Level: logrus.DebugLevel}
	shouldBeNil := g.Fire(&entry)
	if shouldBeNil != nil {
		t.Errorf("Expected nil, on DEBUG, but got %v", shouldBeNil)
		return
	}
}

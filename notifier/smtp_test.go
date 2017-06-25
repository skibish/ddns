package notifier

import (
	"bufio"
	"bytes"
	"errors"
	"net/smtp"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSMTPInit(t *testing.T) {
	_, errConv := initSMTPNotifier("cfg")
	if errConv == nil {
		t.Error("Should be conversion error, but got nothing")
		return
	}

	m := make(map[interface{}]interface{})
	m["user"] = 123
	_, errValue := initSMTPNotifier(m)
	if errValue == nil {
		t.Errorf("Should be error, because value is not a string")
		return
	}

	m = make(map[interface{}]interface{})
	m[123] = "123"
	_, errKey := initSMTPNotifier(m)
	if errKey == nil {
		t.Errorf("Should be error, because key is not a string")
		return
	}

	m = make(map[interface{}]interface{})
	m["user"] = "123"
	s, errUnexpected := initSMTPNotifier(m)
	if errUnexpected != nil {
		t.Error(errUnexpected)
		return
	}

	if s.User != "123" {
		t.Errorf("User should be %q, but got %q", "123", s.User)
		return
	}

}

func TestSMTPFire(t *testing.T) {

	// test, everything is OK
	grabEmailOK := func(string, smtp.Auth, string, []string, []byte) error {
		return nil
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	s := SMTPConfig{
		Host:     "example.com",
		Port:     "1234",
		User:     "yo@example.com",
		Password: "1234",
		To:       "cool@email.address",
		Subject:  "hello, this is test email",
		send:     grabEmailOK,
		errorW:   writer,
	}

	entry := logrus.Entry{Message: "awesome message", Level: logrus.InfoLevel}
	errFire := s.Fire(&entry)
	if errFire != nil {
		t.Errorf("handled an error on SMTP request: %q", errFire.Error())
		return
	}

	// test, got some error from SMTP
	grabEmailBad := func(string, smtp.Auth, string, []string, []byte) error {
		return errors.New("something went wrong")
	}
	s.send = grabEmailBad

	errFire = s.Fire(&entry)
	if errFire == nil {
		t.Error("Should be error, but it's OK")
		return
	}

	// test, DEBUG is ignored
	entry = logrus.Entry{Message: "ignored", Level: logrus.DebugLevel}
	shouldBeNil := s.Fire(&entry)
	if shouldBeNil != nil {
		t.Errorf("Expected nil, on DEBUG, but got %v", shouldBeNil)
		return
	}
}

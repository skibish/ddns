package notifier

import (
	"errors"
	"io"
	"testing"

	"github.com/matryer/is"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func TestSMTPHookNew(t *testing.T) {
	tcases := []struct {
		tname string
		cfg   interface{}
		isErr bool
	}{
		{"ok", map[string]string{"from": "a@a.io", "to": "b@b.io"}, false},
		{"invalid config", "something unexpected", true},
		{"incorrect from", map[string]string{"from": "oh no"}, true},
		{"incorrect to", map[string]string{"from": "oh@no.io", "to": "oh no"}, true},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)

			_, err := newSMTPHook(tc.cfg)

			if tc.isErr {
				if err == nil {
					t.Fail() // should be error
				}
				return
			}
			is.NoErr(err)
		})
	}
}

type mockSender struct {
	send  func() error
	close func() error
}

func (m mockSender) Send(from string, to []string, msg io.WriterTo) error {
	return m.send()
}
func (m mockSender) Close() error {
	return m.close()
}

func TestSMTPHookFire(t *testing.T) {
	is := is.New(t)

	hook, err := newSMTPHook(map[string]interface{}{
		"host":     "smtp.email.server",
		"port":     468,
		"user":     "yo",
		"password": "1234",
		"from":     "yo@yo.co",
		"to":       "bo@bo.co",
		"subject":  "ooomg",
	})
	is.NoErr(err)

	tcases := []struct {
		tname      string
		senderFunc func() (gomail.SendCloser, error)
		isErr      bool
	}{
		{
			tname: "ok",
			senderFunc: func() (gomail.SendCloser, error) {
				m := mockSender{}
				m.send = func() error {
					return nil
				}
				m.close = func() error {
					return nil
				}
				return m, nil
			},
		},
		{
			tname: "dialer failed",
			senderFunc: func() (gomail.SendCloser, error) {
				m := mockSender{}
				m.send = func() error {
					return nil
				}
				m.close = func() error {
					return nil
				}
				return m, errors.New("dialer failed")
			},
			isErr: true,
		},
		{
			tname: "fail send",
			senderFunc: func() (gomail.SendCloser, error) {
				m := mockSender{}
				m.send = func() error {
					return errors.New("failed to send")
				}
				m.close = func() error {
					return nil
				}
				return m, nil
			},
			isErr: true,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)
			hook.senderFunc = tc.senderFunc
			entry := logrus.Entry{Message: "awesome message", Level: logrus.InfoLevel}
			err = hook.Fire(&entry)
			if tc.isErr {
				if err == nil {
					is.Fail() // should be error
				}
				return
			}
			is.NoErr(err)

		})
	}
}

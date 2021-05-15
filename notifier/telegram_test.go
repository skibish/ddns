package notifier

import (
	"net/http"
	"testing"

	"github.com/matryer/is"
	"github.com/sirupsen/logrus"
)

func TestTelegramHookNew(t *testing.T) {
	is := is.New(t)

	if _, err := newTelegramHook("cfg"); err == nil {
		is.Fail() // should be error, but got nothing
	}

	m := make(map[interface{}]interface{})
	m["chat_id"] = "123"
	m["token"] = "tokenv"

	hook, err := newTelegramHook(m)

	is.NoErr(err)
	is.Equal(hook.ChatID, "123")
	is.Equal(hook.Token, "tokenv")
}

func TestTelegramHookFire(t *testing.T) {
	tcases := []struct {
		tname      string
		statusCode int
		isErr      bool
	}{
		{"ok", http.StatusOK, false},
		{"fail do", http.StatusTeapot, true},
		{"fail status", http.StatusTeapot, true},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)

			url, close := httpHelper(t, tc.tname, nil, tc.statusCode)
			defer close()

			hook, err := newTelegramHook(map[string]interface{}{
				"chat_id": "someid",
				"token":   "1234",
			})
			is.NoErr(err)

			hook.host = url
			if tc.tname == "fail do" {
				hook.host = " "
			}

			entry := logrus.Entry{Message: tc.tname, Level: logrus.InfoLevel}
			err = hook.Fire(&entry)

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

package notifier

import (
	"net/http"
	"testing"

	"github.com/matryer/is"
	"github.com/sirupsen/logrus"
)

func TestGotifyHookNew(t *testing.T) {
	is := is.New(t)

	if _, err := newGotifyhook("cfg"); err == nil {
		is.Fail() // should be error, but got nothing
	}

	m := make(map[interface{}]interface{})
	m["app_url"] = 123
	if _, err := newGotifyhook(m); err == nil {
		is.Fail() // should fail because not a valid string
		return
	}

	m = make(map[interface{}]interface{})
	m["app_url"] = "https://gotify.example.com/"
	hook, err := newGotifyhook(m)
	is.NoErr(err)
	is.Equal(hook.AppURL, "https://gotify.example.com")
	is.Equal(hook.Title, "DDNS")

	m["app_url"] = "something bad"
	if _, err := newGotifyhook(m); err == nil {
		is.Fail() // should fail because incorrect url
	}

}

func TestGotifyFire(t *testing.T) {
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

			hook, err := newGotifyhook(map[string]interface{}{
				"app_url": url,
				"token":   "1234",
				"title":   "DDNS",
			})
			is.NoErr(err)

			if tc.tname == "fail do" {
				hook.AppURL = " "
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

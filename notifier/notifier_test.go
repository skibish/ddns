package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func httpHelper(t *testing.T, response string, headers map[string]string, statusCode int) (string, func()) {
	is := is.New(t)
	is.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range headers {
			w.Header().Add(k, v)
		}
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))

	return server.URL, server.Close
}

func TestGetHook(t *testing.T) {
	tcases := []struct {
		tname  string
		config interface{}
		isErr  bool
	}{
		{
			tname: "ok smtp",
			config: map[string]string{
				"type": "smtp",
				"from": "a@a.io",
				"to":   "b@b.io",
			},
		},
		{
			tname: "ok telegram",
			config: map[string]string{
				"type":    "telegram",
				"chat_id": "someid",
				"token":   "1234",
			},
		},
		{
			tname: "ok gotify",
			config: map[string]string{
				"type":    "gotify",
				"app_url": "https://gotify.example.com/",
			},
		},
		{
			tname:  "fail decode type",
			config: map[string]interface{}{"type": 1234},
			isErr:  true,
		},
		{
			tname:  "fail do not exist",
			config: map[string]string{"type": "zzz"},
			isErr:  true,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)

			_, err := GetHook(tc.config)
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

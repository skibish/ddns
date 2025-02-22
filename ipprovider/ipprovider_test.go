package ipprovider

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

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

func TestGetIP(t *testing.T) {
	t.Run("ipv6", func(t *testing.T) {
		is := is.New(t)

		ipp := New(true, 1*time.Second)
		is.True(strings.Contains(ipp.(*IPProvider).providers[0].(*icanhazip).url, "6"))
	})

	tcases := []struct {
		tname    string
		response string
		expected string
		isErr    bool
	}{
		{"ok", `{"ip": "45.45.45.45"}`, "45.45.45.45", false},
		{"fail", "something bad", "45.45.45.45", true},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)

			ipp := New(false, 1*time.Second)

			url, close := httpHelper(t, tc.response, nil, http.StatusOK)
			defer close()

			ipp.(*IPProvider).providers = []ipProvider{
				&ipify{
					c:       &http.Client{},
					url:     url,
					timeout: 1 * time.Second,
				},
			}

			ip, err := ipp.GetIP(context.Background())
			if tc.isErr {
				if err == nil {
					t.Fail() // should be error
				}
				return
			}

			is.Equal(ip, tc.expected)
			is.NoErr(err)
		})
	}
}

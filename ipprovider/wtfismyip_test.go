package ipprovider

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestWtfIsmyip(t *testing.T) {
	t.Run("new ok", func(t *testing.T) {
		is := is.New(t)

		wti := newWtfismyip(10 * time.Second)
		is.Equal(wti.(*wtfIsMyIP).url, "https://ipv4.wtfismyip.com/json")
	})

	t.Run("ipv6 ok", func(t *testing.T) {
		is := is.New(t)

		wti := newWtfismyip(10 * time.Second)
		wti.ForceIPV6()
		is.Equal(wti.(*wtfIsMyIP).url, "https://ipv6.wtfismyip.com/json")
	})

	tcases := []struct {
		tname      string
		response   string
		expected   string
		statusCode int
		isErr      bool
		headers    map[string]string
	}{
		{
			tname:      "ok",
			statusCode: http.StatusOK,
			response:   `{"YourFuckingHostname": "45.45.45.45","YourFuckingIPAddress": "45.45.45.45","YourFuckingISP": "SIA Awesomeness","YourFuckingLocation": "Awesome street","YourFuckingTorExit": "false"}`,
			expected:   "45.45.45.45",
		},
		{
			tname:      "not ok",
			statusCode: http.StatusTeapot,
			response:   `{"ip": "45.45.45.45"}`,
			isErr:      true,
		},
		{
			tname:      "failed to read",
			statusCode: http.StatusOK,
			response:   "something strange",
			isErr:      true,
		},
		{
			tname:      "fail",
			statusCode: http.StatusOK,
			response:   `{"ip": "45.45.45.45"}`,
			isErr:      true,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)

			url, close := httpHelper(t, tc.response, tc.headers, tc.statusCode)
			defer close()

			wti := &wtfIsMyIP{
				c:       &http.Client{},
				url:     url,
				timeout: 10 * time.Second,
			}

			if tc.tname == "fail" {
				wti.url = "localhost:8080"
			}

			ip, err := wti.GetIP(context.Background())
			if tc.isErr {
				if err == nil {
					t.Fail() // should be error
				}
				return
			}
			is.NoErr(err)
			is.Equal(ip, tc.expected)
		})
	}
}

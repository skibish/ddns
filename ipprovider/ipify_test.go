package ipprovider

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestIpify(t *testing.T) {
	t.Run("new ok", func(t *testing.T) {
		is := is.New(t)

		ipf := newIpify(10 * time.Second)
		is.Equal(ipf.(*ipify).url, "https://api.ipify.org/?format=json")
	})

	t.Run("ipv6 ok", func(t *testing.T) {
		is := is.New(t)

		ipf := newIpify(10 * time.Second)
		ipf.ForceIPV6()
		is.Equal(ipf.(*ipify).url, "https://api6.ipify.org/?format=json")

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
			response:   `{"ip": "45.45.45.45"}`,
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

			ipf := &ipify{
				c:       &http.Client{},
				url:     url,
				timeout: 10 * time.Second,
			}

			if tc.tname == "fail" {
				ipf.url = "localhost:8080"
			}

			ip, err := ipf.GetIP(context.Background())
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

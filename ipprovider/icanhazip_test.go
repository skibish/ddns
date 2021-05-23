package ipprovider

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestIcanhazip(t *testing.T) {
	t.Run("new ok", func(t *testing.T) {
		is := is.New(t)

		ifc := newIcanhazip(10 * time.Second)
		is.Equal(ifc.(*icanhazip).url, "https://ipv4.icanhazip.com")
	})

	t.Run("ipv6 ok", func(t *testing.T) {
		is := is.New(t)

		ifc := newIcanhazip(10 * time.Second)
		ifc.ForceIPV6()
		is.Equal(ifc.(*icanhazip).url, "https://ipv6.icanhazip.com")
	})

	tcases := []struct {
		tname      string
		response   string
		statusCode int
		isErr      bool
		headers    map[string]string
	}{
		{
			tname:      "ok",
			statusCode: http.StatusOK,
			response:   "45.45.45.45",
		},
		{
			tname:      "not ok",
			statusCode: http.StatusTeapot,
			response:   "45.45.45.45",
			isErr:      true,
		},
		{
			tname:      "failed to read",
			statusCode: http.StatusOK,
			response:   "45.45.45.45",
			headers: map[string]string{
				"Content-Length": "1",
			},
			isErr: true,
		},
		{
			tname:      "fail",
			statusCode: http.StatusOK,
			response:   "45.45.45.45",
			isErr:      true,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)

			url, close := httpHelper(t, tc.response, tc.headers, tc.statusCode)
			defer close()

			ifc := &icanhazip{
				c:       &http.Client{},
				url:     url,
				timeout: 10 * time.Second,
			}

			if tc.tname == "fail" {
				ifc.url = "localhost:8080"
			}

			ip, err := ifc.GetIP(context.Background())
			if tc.isErr {
				if err == nil {
					t.Fail() // should be error
				}
				return
			}
			is.NoErr(err)
			is.Equal(ip, tc.response)
		})
	}
}

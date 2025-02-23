package do

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func httpHelper(t *testing.T, reqMethod, path, res string, resCode int) (string, func()) {
	is := is.New(t)
	is.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		is.Equal(r.Header.Get("Authorization"), "Bearer amazingtoken")
		is.Equal(r.Method, reqMethod)
		is.Equal(r.URL.Path, path)

		w.WriteHeader(resCode)
		_, _ = w.Write([]byte(res))
	}))

	return server.URL, server.Close
}

func TestList(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	tcases := []struct {
		tname      string
		method     string
		path       string
		doResponse string
		status     int
		isErr      bool
		timeout    time.Duration
	}{
		{
			tname:      "req 200",
			method:     http.MethodGet,
			path:       "/domains/example.com/records",
			doResponse: `{"domain_records":[{"id": 3352895,"type": "A","name": "@","data": "1.2.3.4","priority": null,"port": null,"weight": null}]}`,
			status:     http.StatusOK,
			isErr:      false,
			timeout:    1 * time.Second,
		},
		{
			tname:   "req 500",
			method:  http.MethodGet,
			path:    "/domains/example.com/records",
			status:  http.StatusInternalServerError,
			isErr:   true,
			timeout: 1 * time.Second,
		},
		{
			tname:   "fail to parse DO response",
			method:  http.MethodGet,
			path:    "/domains/example.com/records",
			status:  http.StatusOK,
			isErr:   true,
			timeout: 1 * time.Second,
		},
		{
			tname:      "fail to make a request",
			method:     http.MethodGet,
			path:       "/domains/example.com/records",
			doResponse: `fail`,
			status:     http.StatusOK,
			isErr:      true,
			timeout:    1 * time.Second,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			url, close := httpHelper(t, tc.method, tc.path, tc.doResponse, tc.status)
			defer close()

			d := New("amazingtoken", 1*time.Second)
			d.url = url
			if tc.doResponse == "fail" {
				d.url = "localhost:333"
			}

			recs, err := d.List(context.Background(), "example.com")
			if tc.isErr && err != nil {
				return
			}

			is.NoErr(err)
			is.True(strings.Contains(tc.doResponse, recs[0].Name))
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	tcases := []struct {
		tname      string
		method     string
		path       string
		doResponse string
		status     int
		isErr      bool
		timeout    time.Duration
	}{
		{
			tname:   "req 200",
			method:  http.MethodPost,
			path:    "/domains/example.com/records",
			status:  http.StatusOK,
			isErr:   false,
			timeout: 1 * time.Second,
		},
		{
			tname:   "req 500",
			method:  http.MethodPost,
			path:    "/domains/example.com/records",
			status:  http.StatusInternalServerError,
			isErr:   true,
			timeout: 1 * time.Second,
		},
		{
			tname:   "fail to parse DO response",
			method:  http.MethodPost,
			path:    "/domains/example.com/records",
			status:  http.StatusOK,
			isErr:   true,
			timeout: 1 * time.Second,
		},
		{
			tname:      "fail to make a request",
			method:     http.MethodPost,
			path:       "/domains/example.com/records",
			doResponse: `fail`,
			status:     http.StatusOK,
			isErr:      true,
			timeout:    1 * time.Second,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			url, close := httpHelper(t, tc.method, tc.path, tc.doResponse, tc.status)
			defer close()

			d := New("amazingtoken", 1*time.Second)
			d.url = url
			if tc.doResponse == "fail" {
				d.url = "localhost:333"
			}
			rec := Record{
				Type: "A",
				Name: "@",
				Data: "1.2.3.4",
			}

			err := d.Create(context.Background(), "example.com", rec)
			if tc.isErr && err != nil {
				return
			}

			is.NoErr(err)
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	tcases := []struct {
		tname      string
		method     string
		path       string
		doResponse string
		status     int
		isErr      bool
		timeout    time.Duration
	}{
		{
			tname:   "req 200",
			method:  http.MethodPut,
			path:    "/domains/example.com/records/123",
			status:  http.StatusOK,
			isErr:   false,
			timeout: 1 * time.Second,
		},
		{
			tname:   "req 500",
			method:  http.MethodPut,
			path:    "/domains/example.com/records/123",
			status:  http.StatusInternalServerError,
			isErr:   true,
			timeout: 1 * time.Second,
		},
		{
			tname:   "fail to parse DO response",
			method:  http.MethodPut,
			path:    "/domains/example.com/records/123",
			status:  http.StatusOK,
			isErr:   true,
			timeout: 1 * time.Second,
		},
		{
			tname:      "fail to make a request",
			method:     http.MethodPut,
			path:       "/domains/example.com/records/123",
			doResponse: "fail",
			status:     http.StatusOK,
			isErr:      true,
			timeout:    1 * time.Second,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			url, close := httpHelper(t, tc.method, tc.path, tc.doResponse, tc.status)
			defer close()

			d := New("amazingtoken", 1*time.Second)
			d.url = url
			if tc.doResponse == "fail" {
				d.url = "localhost:333"
			}
			rec := Record{
				ID:   123,
				Type: "A",
				Name: "@",
				Data: "1.2.3.4",
			}

			err := d.Update(context.Background(), "example.com", rec)
			if tc.isErr && err != nil {
				return
			}

			is.NoErr(err)
		})
	}
}

func TestPrepareRequest(t *testing.T) {
	is := is.New(t)

	d := New("amazingtoken", 1*time.Second)
	_, err := d.prepareRequest("12 3", "/path", nil)
	if err == nil {
		is.Fail() // should error because method is incorrect
	}
}

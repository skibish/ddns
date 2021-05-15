package updater

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/skibish/ddns/conf"
	"github.com/skibish/ddns/do"
)

func TestUpdater(t *testing.T) {
	var getIPidx int

	tcases := []struct {
		tname         string
		cfg           *conf.Configuration
		pm            *ProviderMock
		pmGetIPCalls  int
		dm            *DomainsServiceMock
		dmCreateCalls int
		dmUpdateCalls int
		dmListCalls   int
		sleep         time.Duration
	}{
		{
			tname: "ok sync and create",
			cfg: &conf.Configuration{
				Domains: map[string][]do.Record{
					"example.com": {
						{
							Type: "A",
							Name: "ddns",
						},
					},
				},
				CheckPeriod:    3 * time.Second,
				RequestTimeout: 5 * time.Second,
			},
			pm: &ProviderMock{
				GetIPFunc: func(contextMoqParam context.Context) (string, error) {
					return "10.0.5.1", nil
				},
			},
			pmGetIPCalls: 1,
			dm: &DomainsServiceMock{
				CreateFunc: func(contextMoqParam context.Context, s string, record do.Record) error {
					return nil
				},
				ListFunc: func(contextMoqParam context.Context, s string) ([]do.Record, error) {
					return []do.Record{}, nil
				},
			},
			dmCreateCalls: 1,
			dmListCalls:   1,
			sleep:         1 * time.Second,
		},
		{
			tname: "ok sync and update",
			cfg: &conf.Configuration{
				Domains: map[string][]do.Record{
					"example.com": {
						{
							Type: "A",
							Name: "ddns",
						},
						{
							Type: "txt",
							Name: "ddns",
							Data: "updated IP = {{.IP}}, hello, {{.world}}",
						},
					},
				},
				Params: map[string]string{
					"hello": "world",
				},
				CheckPeriod:    1 * time.Second,
				RequestTimeout: 5 * time.Second,
			},
			pm: &ProviderMock{
				GetIPFunc: func(contextMoqParam context.Context) (string, error) {
					ips := []string{"10.0.5.1", "10.0.0.1"}
					ip := ips[getIPidx]
					getIPidx++
					return ip, nil
				},
			},
			pmGetIPCalls: 2,
			dm: &DomainsServiceMock{
				UpdateFunc: func(contextMoqParam context.Context, s string, record do.Record) error {
					return nil
				},
				ListFunc: func(contextMoqParam context.Context, s string) ([]do.Record, error) {
					return []do.Record{
						{
							ID:   123,
							Type: "A",
							Name: "ddns",
						},
						{
							ID:   124,
							Type: "txt",
							Name: "ddns",
						},
					}, nil
				},
			},
			dmUpdateCalls: 4,
			dmListCalls:   2,
			sleep:         1500 * time.Millisecond,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			getIPidx = 0
			is := is.New(t)

			u := New(tc.cfg)
			u.do = tc.dm
			u.ipprovider = tc.pm

			go func() {
				time.Sleep(tc.sleep)
				u.Stop()
			}()

			err := u.Start(context.Background())

			is.NoErr(err)
			is.Equal(len(tc.pm.GetIPCalls()), tc.pmGetIPCalls)
			is.Equal(len(tc.dm.CreateCalls()), tc.dmCreateCalls)
			is.Equal(len(tc.dm.UpdateCalls()), tc.dmUpdateCalls)
			is.Equal(len(tc.dm.ListCalls()), tc.dmListCalls)
		})
	}
}

func TestUpdaterPrepareData(t *testing.T) {
	tcases := []struct {
		tname    string
		input    do.Record
		expected string
		params   map[string]string
		isErr    bool
	}{
		{
			tname:    "ok ip",
			input:    do.Record{},
			params:   make(map[string]string),
			expected: "10.0.0.1",
		},
		{
			tname: "ok template ip",
			input: do.Record{
				Type: "TXT",
				Data: "Hello {{.IP}}",
			},
			params:   make(map[string]string),
			expected: "Hello 10.0.0.1",
		},
		{
			tname: "ok template with params",
			input: do.Record{
				Type: "TXT",
				Data: "Hello {{.IP}}, {{.myprop}}",
			},
			params: map[string]string{
				"myprop": "hello",
			},
			expected: "Hello 10.0.0.1, hello",
		},
		{
			tname: "failed to parse the template",
			input: do.Record{
				Type: "TXT",
				Data: "Hello {{.IP}}, {{.myprop}",
			},
			params: make(map[string]string),
			isErr:  true,
		},
	}

	u := &Updater{
		ip: "10.0.0.1",
	}
	for _, tc := range tcases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)

			v, err := u.prepareData(tc.input, tc.params)

			if tc.isErr {
				if err == nil {
					is.Fail()
				}
				return
			}

			is.Equal(v, tc.expected)
			is.NoErr(err)
		})
	}

}

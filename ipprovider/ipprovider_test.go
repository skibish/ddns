package ipprovider

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestGetIP(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	log.SetOutput(bw)

	handler := func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`{
    "YourFuckingIPAddress": "45.45.45.45",
    "ip": "45.45.45.45"
}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ch := &http.Client{}

	providers = append(providers,
		&Ifconfig{c: ch, url: "http://127.0.0.1:1234"},
		&Ipify{c: ch, url: server.URL})

	ip := GetIP()
	if ip != "45.45.45.45" {
		t.Error("aaa")
		return
	}

	errFlush := bw.Flush()
	if errFlush != nil {
		t.Errorf("Error flushing log message: %q", errFlush.Error())
		return
	}

	if !strings.Contains(b.String(), "level=warning") {
		t.Error("There should be level=warning in log")
		return
	}

}

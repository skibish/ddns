package ifconfig

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIfConfigNew(t *testing.T) {
	expectedURL := "https://v4.ifconfig.co/json"
	ifc := New(&http.Client{}, false)
	ifcOriginal := ifc.(*ifconfig)
	if ifcOriginal.url != expectedURL {
		t.Errorf("URL of ifconfig should be %q, but got %q", expectedURL, ifcOriginal.url)
		return
	}

	expectedv6URL := "https://v6.ifconfig.co/json"
	ifc = New(&http.Client{}, true)
	ifcOriginal = ifc.(*ifconfig)
	if ifcOriginal.url != expectedv6URL {
		t.Errorf("URL of ifconfig should be %q, but got %q", expectedv6URL, ifcOriginal.url)
		return
	}
}

func TestIfConfigSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`{"city": "Unknown","country": "Latvia","ip": "45.45.45.45","ip_decimal": 1424881195}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ifc := &ifconfig{
		c:   &http.Client{},
		url: server.URL,
	}

	ip, errGet := ifc.GetIP()
	if errGet != nil {
		t.Errorf("Got error: %s", errGet.Error())
		return
	}

	if ip != "45.45.45.45" {
		t.Errorf("Incorrect IP value. Got %q, but should be %q", ip, "45.45.45.45")
		return
	}
}

func TestIfConfigNotSuccessCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(429)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ifc := &ifconfig{
		c:   &http.Client{},
		url: server.URL,
	}

	_, errGet := ifc.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if errGet.Error() != "ifconfig: Status code is not in success range: 429" {
		t.Error("Error was, but not about status code")
		return
	}
}

func TestIfConfigFailedDecode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`something wrong`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ifc := &ifconfig{
		c:   &http.Client{},
		url: server.URL,
	}

	_, errGet := ifc.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if errGet.Error() != "ifconfig: invalid character 's' looking for beginning of value" {
		t.Error("Error was, but not related to parsing")
		return
	}
}

func TestIfConfigFailedOnGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`something wrong`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ifc := &ifconfig{
		c:   &http.Client{},
		url: "http://127.0.0.1:1234",
	}

	_, errGet := ifc.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if errGet.Error() != "ifconfig: Get http://127.0.0.1:1234: dial tcp 127.0.0.1:1234: connect: connection refused" {
		t.Errorf("Error was, but not related to the request fail: %v", errGet.Error())
		return
	}
}

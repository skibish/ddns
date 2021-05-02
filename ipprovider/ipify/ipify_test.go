package ipify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIpifyNew(t *testing.T) {
	expectedURL := "https://api.ipify.org/?format=json"
	ipf := New(&http.Client{})
	ipfOriginal := ipf.(*ipify)
	if ipfOriginal.url != expectedURL {
		t.Errorf("URL of ipfonfig should be %q, but got %q", expectedURL, ipfOriginal.url)
		return
	}
}

func TestForceIPV6(t *testing.T) {
	expectedv6URL := "https://api6.ipify.org/?format=json"
	ipf := New(&http.Client{})
	ipf.ForceIPV6()
	ipfOriginal := ipf.(*ipify)
	if ipfOriginal.url != expectedv6URL {
		t.Errorf("URL of ipfonfig should be %q, but got %q", expectedv6URL, ipfOriginal.url)
		return
	}
}

func TestIpifySuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`{"ip": "45.45.45.45"}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ipf := &ipify{
		c:   &http.Client{},
		url: server.URL,
	}

	ip, errGet := ipf.GetIP()
	if errGet != nil {
		t.Errorf("Got error: %s", errGet.Error())
		return
	}

	if ip != "45.45.45.45" {
		t.Errorf("Incorrect IP value. Got %q, but should be %q", ip, "45.45.45.45")
		return
	}
}

func TestIpifyNotSuccessCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(429)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ipf := &ipify{
		c:   &http.Client{},
		url: server.URL,
	}

	_, errGet := ipf.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if errGet.Error() != "ipify: Status code is not in success range: 429" {
		t.Error("Error was, but not about status code")
		return
	}
}

func TestIpifyFailedDecode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`something wrong`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ipf := &ipify{
		c:   &http.Client{},
		url: server.URL,
	}

	_, errGet := ipf.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if errGet.Error() != "ipify: invalid character 's' looking for beginning of value" {
		t.Error("Error was, but not related to parsing")
		return
	}
}

func TestIpifyFailedOnGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`something wrong`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ipf := &ipify{
		c:   &http.Client{},
		url: "http://127.0.0.1:1234",
	}

	_, errGet := ipf.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}
}

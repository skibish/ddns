package wtfismyip

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWtfIsMyIPNew(t *testing.T) {
	expectedURL := "https://ipv4.wtfismyip.com/json"
	wti := New(&http.Client{})
	wtiOriginal := wti.(*wtfIsMyIP)
	if wtiOriginal.url != expectedURL {
		t.Errorf("URL of wtionfig should be %q, but got %q", expectedURL, wtiOriginal.url)
		return
	}
}

func TestForceIPV6(t *testing.T) {
	expectedURL := "https://ipv6.wtfismyip.com/json"
	wti := New(&http.Client{})
	wti.ForceIPV6()
	wtiOriginal := wti.(*wtfIsMyIP)
	if wtiOriginal.url != expectedURL {
		t.Errorf("URL of wtionfig should be %q, but got %q", expectedURL, wtiOriginal.url)
		return
	}
}

func TestWtfIsMyIPSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`{
    "YourFuckingHostname": "45.45.45.45",
    "YourFuckingIPAddress": "45.45.45.45",
    "YourFuckingISP": "SIA Awesomeness",
    "YourFuckingLocation": "Awesome street",
    "YourFuckingTorExit": "false"
}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	wti := &wtfIsMyIP{
		c:   &http.Client{},
		url: server.URL,
	}

	ip, errGet := wti.GetIP()
	if errGet != nil {
		t.Errorf("Got error: %s", errGet.Error())
		return
	}

	if ip != "45.45.45.45" {
		t.Errorf("Incorrect IP value. Got %q, but should be %q", ip, "45.45.45.45")
		return
	}
}

func TestWtfIsMyIPNotSuccessCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(429)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	wti := &wtfIsMyIP{
		c:   &http.Client{},
		url: server.URL,
	}

	_, errGet := wti.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if errGet.Error() != "wtfismyip: Status code is not in success range: 429" {
		t.Error("Error was, but not about status code")
		return
	}
}

func TestWtfIsMyIPFailedDecode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`something wrong`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	wti := &wtfIsMyIP{
		c:   &http.Client{},
		url: server.URL,
	}

	_, errGet := wti.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if errGet.Error() != "wtfismyip: invalid character 's' looking for beginning of value" {
		t.Error("Error was, but not related to parsing")
		return
	}
}

func TestWtfIsMyIPFailedOnGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`something wrong`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	wti := &wtfIsMyIP{
		c:   &http.Client{},
		url: "http://127.0.0.1:1234",
	}

	_, errGet := wti.GetIP()
	if errGet == nil {
		t.Errorf("Should be error, but is success")
		return
	}

	if !isMatchingErrorMessage(errGet.Error(), "wtfismyip", "connection refused") {
		t.Error("Error was, but not related to the request fail")
		return
	}
}

func isMatchingErrorMessage(message string, prefix, suffix string) bool {
	return strings.HasPrefix(message, prefix) && strings.HasSuffix(message, suffix)
}

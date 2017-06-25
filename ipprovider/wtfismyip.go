package ipprovider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skibish/ddns/misc"
)

// WtfIsMyIP is an abstraction to get IP
type WtfIsMyIP struct {
	c   *http.Client
	url string
}

// wtfIsMyIPResponse is a response
type wtfIsMyIPResponse struct {
	IP string `json:"YourFuckingIPAddress"`
}

// newWtfIsMyIP return WtfIsMyIP
func newWtfIsMyIP(c *http.Client) *WtfIsMyIP {
	return &WtfIsMyIP{
		c:   c,
		url: "https://wtfismyip.com/json",
	}
}

// GetIP get IP
func (i *WtfIsMyIP) GetIP() (string, error) {
	resp, errGet := i.c.Get(i.url)
	if errGet != nil {
		return "", fmt.Errorf("%s: %s", wtfismyipName, errGet.Error())
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("%s: Status code is not in success range: %d", wtfismyipName, resp.StatusCode)
	}

	var r wtfIsMyIPResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&r)
	if errDecode != nil {
		return "", fmt.Errorf("%s: %s", wtfismyipName, errDecode.Error())
	}

	return r.IP, nil
}

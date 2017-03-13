package ipprovider

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		return "", errGet
	}

	defer resp.Body.Close()

	if !success(resp.StatusCode) {
		return "", fmt.Errorf("Status code is not in success range: %d", resp.StatusCode)
	}

	var r wtfIsMyIPResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&r)
	if errDecode != nil {
		return "", errDecode
	}

	return r.IP, nil
}

package ipprovider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skibish/ddns/misc"
)

// Ipify is an abstraction to get IP
type Ipify struct {
	c   *http.Client
	url string
}

type ipifyResponse struct {
	IP string `json:"ip"`
}

// newIpify return Ipify
func newIpify(c *http.Client) *Ipify {
	return &Ipify{
		c:   c,
		url: "https://api.ipify.org/?format=json",
	}
}

// GetIP get ip
func (i *Ipify) GetIP() (string, error) {
	resp, errGet := i.c.Get(i.url)
	if errGet != nil {
		return "", errGet
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("Status code is not in success range: %d", resp.StatusCode)
	}

	var r ipifyResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&r)
	if errDecode != nil {
		return "", errDecode
	}

	return r.IP, nil
}

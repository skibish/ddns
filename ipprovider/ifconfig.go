package ipprovider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skibish/ddns/misc"
)

// Ifconfig is an abstraction to get IP
type Ifconfig struct {
	c   *http.Client
	url string
}

// IfconfigResponse is a response
type ifconfigResponse struct {
	IP string `json:"ip"`
}

// newIfConfig return Ifconfig
func newIfConfig(c *http.Client) *Ifconfig {
	return &Ifconfig{
		c:   c,
		url: "https://ifconfig.co/json",
	}
}

// GetIP get ip
func (i *Ifconfig) GetIP() (string, error) {
	resp, errGet := i.c.Get(i.url)
	if errGet != nil {
		return "", errGet
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("Status code is not in success range: %d", resp.StatusCode)
	}

	var r ifconfigResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&r)
	if errDecode != nil {
		return "", errDecode
	}

	return r.IP, nil
}

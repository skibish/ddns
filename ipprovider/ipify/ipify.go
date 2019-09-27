package ipify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skibish/ddns/ipprovider"

	"github.com/skibish/ddns/misc"
)

var providerName = "ipify"

// ipify is an abstraction to get IP
type ipify struct {
	c   *http.Client
	url string
}

type ipifyResponse struct {
	IP string `json:"ip"`
}

// New return Ipify
func New(c *http.Client) ipprovider.Provider {
	return &ipify{
		c:   c,
		url: "https://api.ipify.org/?format=json",
	}
}

// ForceIPV6 .
func (i *ipify) ForceIPV6() {
	i.url = "https://api6.ipify.org/?format=json"
}

// GetIP get ip
func (i *ipify) GetIP() (string, error) {
	resp, errGet := i.c.Get(i.url)
	if errGet != nil {
		return "", fmt.Errorf("%s: %s", providerName, errGet.Error())
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("%s: Status code is not in success range: %d", providerName, resp.StatusCode)
	}

	var r ipifyResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&r)
	if errDecode != nil {
		return "", fmt.Errorf("%s: %s", providerName, errDecode.Error())
	}

	return r.IP, nil
}

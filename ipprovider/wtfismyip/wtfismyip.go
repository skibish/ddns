package wtfismyip

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skibish/ddns/ipprovider"

	"github.com/skibish/ddns/misc"
)

var providerName = "wtfismyip"

// wtfIsMyIP is an abstraction to get IP
type wtfIsMyIP struct {
	c   *http.Client
	url string
}

// wtfIsMyIPResponse is a response
type wtfIsMyIPResponse struct {
	IP string `json:"YourFuckingIPAddress"`
}

// New return wtfIsMyIP
func New(c *http.Client) ipprovider.Provider {
	return &wtfIsMyIP{
		c:   c,
		url: "https://ipv4.wtfismyip.com/json",
	}
}

// ForceIPV6 .
func (i *wtfIsMyIP) ForceIPV6() {
	i.url = "https://ipv6.wtfismyip.com/json"
}

// GetIP get IP
func (i *wtfIsMyIP) GetIP() (string, error) {
	resp, errGet := i.c.Get(i.url)
	if errGet != nil {
		return "", fmt.Errorf("%s: %s", providerName, errGet.Error())
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("%s: Status code is not in success range: %d", providerName, resp.StatusCode)
	}

	var r wtfIsMyIPResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&r)
	if errDecode != nil {
		return "", fmt.Errorf("%s: %s", providerName, errDecode.Error())
	}

	return r.IP, nil
}

package ifconfig

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skibish/ddns/ipprovider"
	"github.com/skibish/ddns/misc"
)

var providerName = "ifconfig"

// ifconfig is an abstraction to get IP
type ifconfig struct {
	c   *http.Client
	url string
}

// ifconfigResponse is a response
type ifconfigResponse struct {
	IP string `json:"ip"`
}

// New return ifconfig
func New(c *http.Client) ipprovider.Provider {
	return &ifconfig{
		c:   c,
		url: "https://v4.ifconfig.co/json",
	}
}

// ForceIPV6 .
func (i *ifconfig) ForceIPV6() {
	i.url = "https://v6.ifconfig.co/json"
}

// GetIP get ip
func (i *ifconfig) GetIP() (string, error) {
	resp, errGet := i.c.Get(i.url)
	if errGet != nil {
		return "", fmt.Errorf("%s: %s", providerName, errGet.Error())
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("%s: Status code is not in success range: %d", providerName, resp.StatusCode)
	}

	var r ifconfigResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&r)
	if errDecode != nil {
		return "", fmt.Errorf("%s: %s", providerName, errDecode.Error())
	}

	return r.IP, nil
}

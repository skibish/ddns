package icanhazip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/skibish/ddns/ipprovider"
	"github.com/skibish/ddns/misc"
)

var providerName = "icanhazip"

type icanhazip struct {
	c   *http.Client
	url string
}

// New return icanhazip
func New(c *http.Client) ipprovider.Provider {
	return &icanhazip{
		c:   c,
		url: "https://ipv4.icanhazip.com",
	}
}

// ForceIPV6
func (i *icanhazip) ForceIPV6() {
	i.url = "https://ipv6.icanhazip.com"
}

func (i *icanhazip) GetIP() (string, error) {
	resp, errGet := i.c.Get(i.url)
	if errGet != nil {
		return "", fmt.Errorf("%s: %s", providerName, errGet.Error())
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("%s: Status code is not in success range: %d", providerName, resp.StatusCode)
	}

	b, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return "", fmt.Errorf("%s: Failed to read body of response: %d", providerName, resp.StatusCode)
	}

	return strings.TrimSpace(string(b)), nil
}

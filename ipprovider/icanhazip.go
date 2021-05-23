package ipprovider

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/skibish/ddns/misc"
)

type icanhazip struct {
	c       *http.Client
	url     string
	timeout time.Duration
}

// New return icanhazip
func newIcanhazip(timeout time.Duration) ipProvider {
	return &icanhazip{
		c:       &http.Client{},
		url:     "https://ipv4.icanhazip.com",
		timeout: timeout,
	}
}

// ForceIPV6
func (i *icanhazip) ForceIPV6() {
	i.url = "https://ipv6.icanhazip.com"
}

func (i *icanhazip) GetIP(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, i.url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create a request: %w", err)
	}

	resp, err := i.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("status code is not in success range: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the body of the response: %d", resp.StatusCode)
	}

	return strings.TrimSpace(string(b)), nil
}

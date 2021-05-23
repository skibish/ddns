package ipprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skibish/ddns/misc"
)

// wtfIsMyIP is an abstraction to get IP
type wtfIsMyIP struct {
	c       *http.Client
	url     string
	timeout time.Duration
}

// wtfIsMyIPResponse is a response
type wtfIsMyIPResponse struct {
	IP string `json:"YourFuckingIPAddress"`
}

// New return wtfIsMyIP
func newWtfismyip(timeout time.Duration) ipProvider {
	return &wtfIsMyIP{
		c:       &http.Client{},
		url:     "https://ipv4.wtfismyip.com/json",
		timeout: timeout,
	}
}

// ForceIPV6 .
func (i *wtfIsMyIP) ForceIPV6() {
	i.url = "https://ipv6.wtfismyip.com/json"
}

// GetIP get IP
func (i *wtfIsMyIP) GetIP(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, i.url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create a request: %w", err)
	}

	resp, err := i.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to do a request: %w", err)
	}

	defer resp.Body.Close()

	if !misc.Success(resp.StatusCode) {
		return "", fmt.Errorf("status code is not in success range: %d", resp.StatusCode)
	}

	var r wtfIsMyIPResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("failed to decode the response: %w", err)
	}

	return r.IP, nil
}

package ipprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skibish/ddns/misc"
)

// ipify is an abstraction to get IP
type ipify struct {
	c       *http.Client
	url     string
	timeout time.Duration
}

type ipifyResponse struct {
	IP string `json:"ip"`
}

// New return Ipify
func newIpify(timeout time.Duration) ipProvider {
	return &ipify{
		c:       &http.Client{},
		url:     "https://api.ipify.org/?format=json",
		timeout: timeout,
	}
}

// ForceIPV6 .
func (i *ipify) ForceIPV6() {
	i.url = "https://api6.ipify.org/?format=json"
}

// GetIP get ip
func (i *ipify) GetIP(ctx context.Context) (string, error) {
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

	var r ipifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("failed to decode the response: %w", err)
	}

	return r.IP, nil
}

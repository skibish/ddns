package ipprovider

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

// Provider is an interface that should
// be implemented by all IP providers.
type Provider interface {
	GetIP(context.Context) (string, error)
}

type ipProvider interface {
	Provider
	ForceIPV6()
}

// IPProvider struct is IP provider service.
type IPProvider struct {
	providers []ipProvider
}

// New return new IPProvider instance.
func New(ipv6 bool, timeout time.Duration) Provider {
	providers := []ipProvider{
		newIcanhazip(timeout),
		newWtfismyip(timeout),
		newIpify(timeout),
	}

	if ipv6 {
		for _, p := range providers {
			p.ForceIPV6()
		}
	}

	return &IPProvider{
		providers: providers,
	}
}

// GetIP return IP from the first successful source.
func (i *IPProvider) GetIP(ctx context.Context) (string, error) {
	for _, p := range i.providers {
		ip, err := p.GetIP(ctx)
		if err != nil {
			log.Warn(err)
		}
		if ip != "" {
			return ip, nil
		}
	}

	return "", errors.New("failed to get ip from providers")
}

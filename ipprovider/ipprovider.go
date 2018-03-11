package ipprovider

import (
	log "github.com/sirupsen/logrus"
)

// Provider is an interface that should
// be implemented by all IP providers.
type Provider interface {
	GetIP() (string, error)
	ForceIPV6()
}

// IPProvider struct
type IPProvider struct {
	providers []Provider
}

// New return new IPProvider instance
func New() *IPProvider {
	return &IPProvider{}
}

// Register registers IP providers
func (i *IPProvider) Register(p ...Provider) {
	i.providers = append(i.providers, p...)
}

// GetIP return ip from first successful source
func (i *IPProvider) GetIP() (ip string) {
	for _, p := range i.providers {
		var errGet error

		ip, errGet = p.GetIP()
		if errGet != nil {
			log.Warn(errGet.Error())
		}
		if ip != "" {
			break
		}
	}

	return ip
}

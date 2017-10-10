package ipprovider

import (
	log "github.com/sirupsen/logrus"
)

// list of providers
var providers []Provider

// Provider is an interface that should
// be implemented by all IP providers.
type Provider interface {
	GetIP() (string, error)
}

// Register registers IP providers
func Register(p ...Provider) {
	providers = append(providers, p...)
}

// GetIP return ip from first successful source
func GetIP() (ip string) {
	for _, p := range providers {
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

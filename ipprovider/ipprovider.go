package ipprovider

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

var providers []IPProvider

// IPProvider is an interface that should
// be implemented by all IP providers.
type IPProvider interface {
	GetIP() (string, error)
}

// Register registers IP providers
func Register(c *http.Client) {
	providers = append(providers,
		newIfConfig(c), newIpify(c), newWtfIsMyIP(c))
}

// GetIP retur ip from first successful source
func GetIP() (ip string) {
	for _, v := range providers {
		var errGet error
		ip, errGet = v.GetIP()
		if errGet != nil {
			log.Warn(errGet.Error())
		}
		if ip != "" {
			break
		}
	}

	return ip
}

func success(status int) bool {
	return status >= 200 && status < 300
}

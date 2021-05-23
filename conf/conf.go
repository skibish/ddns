package conf

import (
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skibish/ddns/do"
	"github.com/spf13/viper"
)

// Configuration is a structure which holds DDNS configuration.
type Configuration struct {
	Token          string
	IPv6           bool
	CheckPeriod    time.Duration
	RequestTimeout time.Duration
	Domains        map[string][]do.Record
	Notifications  []map[string]interface{}
	Params         map[string]string
}

// valid checks that provided configuration is valid
func (c *Configuration) valid() error {
	if c.Token == "" {
		return errors.New("token can't be empty")
	}

	if len(c.Domains) == 0 {
		return errors.New("domains can't be empty")
	}

	if len(c.Domains) > 0 {
		for domain, records := range c.Domains {
			if len(records) == 0 {
				return fmt.Errorf("records can't be empty for %s", domain)
			}
		}
	}

	return nil
}

// NewConfiguration read configuration file
// and return *Configuration
func NewConfiguration(path string) (*Configuration, error) {
	v := viper.NewWithOptions(
		viper.KeyDelimiter("::"),
		viper.EnvKeyReplacer(strings.NewReplacer("::", "_")),
	)

	v.SetConfigName("ddns")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")

	v.SetDefault("CheckPeriod", 5*time.Minute)
	v.SetDefault("RequestTimeout", 10*time.Second)
	v.SetDefault("IPv6", false)

	if path != "" {
		v.SetConfigFile(path)
	}

	v.SetEnvPrefix("ddns")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("configuration file not found")
		}

		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}
	log.Debugf("using the following configuration file: %s", v.ConfigFileUsed())

	var cf Configuration
	if err := v.Unmarshal(&cf); err != nil {
		return nil, err
	}

	errValid := cf.valid()
	if errValid != nil {
		return nil, errValid
	}

	if cf.Params == nil {
		cf.Params = make(map[string]string)
	}

	return &cf, nil
}

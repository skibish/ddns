package notifier

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type HookType struct {
	Type string
}

// GetHook returns initialized notifier as hook for Logrus
func GetHook(cfg interface{}) (logrus.Hook, error) {
	var ht HookType
	if err := mapstructure.Decode(cfg, &ht); err != nil {
		return nil, err
	}

	switch strings.ToLower(ht.Type) {
	case "smtp":
		return newSMTPHook(cfg)
	case "telegram":
		return newTelegramHook(cfg)
	case "gotify":
		return newGotifyhook(cfg)
	default:
		return nil, fmt.Errorf("notifier %s does not exists", ht.Type)
	}
}

func AllowedLevels() []logrus.Level {
	// ignore DEBUG messages
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

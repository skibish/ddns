package notifier

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// GetHook returns initialized notifier as hook for Logrus
func GetHook(name string, cfg interface{}) (logrus.Hook, error) {
	switch name {
	case "smtp":
		return initSMTPNotifier(cfg)
	default:
		return nil, fmt.Errorf("Notifier %q does not exist", name)
	}
}

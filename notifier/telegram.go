package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/skibish/ddns/misc"
)

// telegramHook is a structure for Telegram notifications configuration
type telegramHook struct {
	Token  string
	ChatID string `mapstructure:"chat_id"`
	host   string
	c      *http.Client
}

func newTelegramHook(cfg interface{}) (*telegramHook, error) {
	var hook telegramHook
	if err := mapstructure.Decode(cfg, &hook); err != nil {
		return nil, err
	}

	hook.host = "https://api.telegram.org"
	hook.c = &http.Client{}

	return &hook, nil
}

func (hook *telegramHook) send(msg string) error {
	form := &url.Values{}
	form.Add("chat_id", hook.ChatID)
	form.Add("text", msg)

	url := fmt.Sprintf("%s/bot%s/sendMessage", hook.host, hook.Token)
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create a request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := hook.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do a request: %w", err)
	}
	defer res.Body.Close()

	if !misc.Success(res.StatusCode) {
		return fmt.Errorf("status code is not in a success range: %w", err)
	}

	return nil
}

// Fire fires hook
func (hook *telegramHook) Fire(entry *logrus.Entry) error {
	if err := hook.send(entry.Message); err != nil {
		return fmt.Errorf("failed to send a message to the chat: %w", err)
	}

	return nil
}

// Levels return array of levels
func (hook *telegramHook) Levels() []logrus.Level {
	return allowedLevels()
}

package notifier

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/skibish/ddns/misc"
)

// gotifyHook is a structure for Gotify notifications configuration
type gotifyHook struct {
	AppURL   string `mapstructure:"app_url"`
	AppToken string `mapstructure:"app_token"`
	Title    string
	c        *http.Client
}

func newGotifyhook(cfg interface{}) (*gotifyHook, error) {
	var hook gotifyHook
	if err := mapstructure.Decode(cfg, &hook); err != nil {
		return nil, err
	}

	if hook.Title == "" {
		hook.Title = "DDNS"
	}

	if !isValidURL(hook.AppURL) {
		return nil, errors.New("app_url is not a valid url")
	}

	hook.AppURL = strings.TrimSuffix(hook.AppURL, "/")
	hook.c = &http.Client{}

	return &hook, nil
}

func (hook *gotifyHook) send(msg string) error {
	form := &url.Values{}
	form.Add("title", hook.Title)
	form.Add("message", msg)

	url := fmt.Sprintf("%s/message?token=%s", hook.AppURL, hook.AppToken)
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
func (hook *gotifyHook) Fire(entry *logrus.Entry) error {
	if err := hook.send(entry.Message); err != nil {
		return fmt.Errorf("failed to send a message to app %w", err)
	}

	return nil
}

// Levels return array of levels
func (hook *gotifyHook) Levels() []logrus.Level {
	return AllowedLevels()
}

// Source: https://golangcode.com/how-to-check-if-a-string-is-a-url/
// isValidURL tests a string to determine if it is a well-structured url or not.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

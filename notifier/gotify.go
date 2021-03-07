package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// GotifyConfig is a structure for Gotify notifications configuration
type GotifyConfig struct {
	AppURL   string `json:"app_url"`
	AppToken string `json:"app_token"`
	send     func(string) error
	errorW   io.Writer
}

func initGotifyNotifier(cfg interface{}) (*GotifyConfig, error) {
	// because in YAML we can have keys of complex type, they are usually of type
	// map[interface{}]interface{}. In case for this hook we are interesting to convert
	// it to map[string]string.
	originalCfg, ok := cfg.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("not converted passed configuration")
	}
	m2 := make(map[string]interface{})

	for key, value := range originalCfg {
		switch key := key.(type) {
		case string:
			m2[key] = value
		default:
			return nil, errors.New("all keys should be strings in YAML")
		}
	}

	// here we mashalling-unmarshalling to fill structure with correct values
	b, errMarshal := json.Marshal(m2)
	if errMarshal != nil {
		return nil, errMarshal
	}

	var g GotifyConfig
	errUnm := json.Unmarshal(b, &g)
	if errUnm != nil {
		return nil, errUnm
	}

	if !isValidURL(g.AppURL) {
		return nil, errors.New("app_url is not a valid URL")
	}

	g.AppURL = strings.TrimSuffix(g.AppURL, "/")

	g.errorW = os.Stdout
	// function that send notification to gotify
	g.send = func(msg string) error {
		u := fmt.Sprintf("%s/message?token=%s", g.AppURL, g.AppToken)
		data := url.Values{
			"title":   {"DDNS"},
			"message": {msg},
		}

		_, err := http.PostForm(u, data)
		if err != nil {
			return err
		}

		return nil
	}

	return &g, nil
}

// Fire fires hook
func (g *GotifyConfig) Fire(entry *logrus.Entry) error {
	// ignoring recording of events on DEBUG
	if entry.Level == logrus.DebugLevel {
		return nil
	}

	errSend := g.send(entry.Message)
	if errSend != nil {
		fmt.Fprintf(g.errorW, "unable to send message to app: %v\n", errSend)
		return errSend
	}

	return nil
}

// Levels return array of levels
func (g *GotifyConfig) Levels() []logrus.Level {
	return logrus.AllLevels
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

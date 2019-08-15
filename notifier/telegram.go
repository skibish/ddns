package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
)

// TelegramConfig is a structure for Telegram notifications configuration
type TelegramConfig struct {
	Token  string `json:"token"`
	ChatID string `json:"chat_id"`
	send   func(string) error
	errorW io.Writer
}

func initTelegramNotifier(cfg interface{}) (*TelegramConfig, error) {
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

	var tg TelegramConfig
	errUnm := json.Unmarshal(b, &tg)
	if errUnm != nil {
		return nil, errUnm
	}

	tg.errorW = os.Stdout
	// function that send notification to telegram
	tg.send = func(msg string) error {
		u := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tg.Token)
		data := url.Values{
			"chat_id": {tg.ChatID},
			"text":    {msg},
		}

		_, err := http.PostForm(u, data)
		if err != nil {
			return err
		}

		return nil
	}

	return &tg, nil
}

// Fire fires hook
func (tg *TelegramConfig) Fire(entry *logrus.Entry) error {
	// ignoring recording of events on DEBUG
	if entry.Level == logrus.DebugLevel {
		return nil
	}

	errSend := tg.send(entry.Message)
	if errSend != nil {
		fmt.Fprintf(tg.errorW, "unable to send message to chat: %v\n", errSend)
		return errSend
	}

	return nil
}

// Levels return array of levels
func (tg *TelegramConfig) Levels() []logrus.Level {
	return logrus.AllLevels
}

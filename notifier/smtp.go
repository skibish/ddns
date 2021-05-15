package notifier

import (
	"fmt"
	"net/mail"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// smtpHook is a structure for smtpHook configuration
type smtpHook struct {
	Host       string
	Port       int
	User       string
	Password   string
	From       string
	To         string
	Subject    string
	senderFunc func() (gomail.SendCloser, error)
}

// newSMTPHook initializes SMTPConfig structure.
func newSMTPHook(cfg interface{}) (*smtpHook, error) {
	var hook smtpHook
	if err := mapstructure.Decode(cfg, &hook); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %v", err)
	}

	if _, err := mail.ParseAddress(hook.From); err != nil {
		return nil, fmt.Errorf("failed to parse from address: %w", err)
	}

	if _, err := mail.ParseAddress(hook.To); err != nil {
		return nil, fmt.Errorf("failed to parse to address: %w", err)
	}

	hook.senderFunc = func() (gomail.SendCloser, error) {
		d := gomail.NewDialer(hook.Host, hook.Port, hook.User, hook.Password)
		s, err := d.Dial()
		if err != nil {
			return nil, err
		}

		return s, nil
	}

	return &hook, nil
}

// Fire fires hook
func (hook *smtpHook) Fire(entry *logrus.Entry) error {
	m := gomail.NewMessage()
	m.SetHeader("From", hook.From)
	m.SetHeader("To", hook.To)
	m.SetHeader("Subject", hook.Subject)
	m.SetBody("text/html", entry.Message)
	s, err := hook.senderFunc()
	if err != nil {
		return err
	}
	defer s.Close()

	if err := gomail.Send(s, m); err != nil {
		return err
	}

	return nil
}

// Levels return array of levels
func (hook *smtpHook) Levels() []logrus.Level {
	return AllowedLevels()
}

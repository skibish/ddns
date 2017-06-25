package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/smtp"
	"os"

	"github.com/sirupsen/logrus"
)

var smtpTmpl = "From: %s\nTo: %s\nSubject: %s\n\n%s"

// SMTPConfig is a structure for SMTPConfig configuration
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	send     func(string, smtp.Auth, string, []string, []byte) error
	errorW   io.Writer
}

// initSMTPNotifier initializes SMTPConfig structure.
func initSMTPNotifier(cfg interface{}) (*SMTPConfig, error) {

	// because in YAML we can have keys of complex type, they are usually of type
	// map[interface{}]interface{}. In case for this hook we are interesting to convert
	// it to map[string]string.
	originalCfg, ok := cfg.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("Not converted passed configuration")
	}
	m2 := make(map[string]string)

	for key, value := range originalCfg {
		switch key := key.(type) {
		case string:
			switch value := value.(type) {
			case string:
				m2[key] = value
			default:
				return nil, errors.New("All values should be strings in YAML")
			}
		default:
			return nil, errors.New("All keys should be strings in YAML")
		}
	}

	// here we mashalling-unmarshalling to fill structure with correct values
	b, errMarshal := json.Marshal(m2)
	if errMarshal != nil {
		return nil, errMarshal
	}

	var s SMTPConfig
	errUnm := json.Unmarshal(b, &s)
	if errUnm != nil {
		return nil, errUnm
	}

	// we are not using smtp.SendMail dirrectly, because we want to test
	// Fire() method
	s.send = smtp.SendMail
	s.errorW = os.Stdout

	return &s, nil
}

// Fire fires hook
func (s *SMTPConfig) Fire(entry *logrus.Entry) error {
	// ignoring recoring of events on DEBUG
	if entry.Level == logrus.DebugLevel {
		return nil
	}

	msg := fmt.Sprintf(smtpTmpl, s.User, s.To, s.Subject, entry.Message)

	auth := smtp.PlainAuth("", s.User, s.Password, s.Host)
	errSend := s.send(s.Host+":"+s.Port, auth, s.User, []string{s.To}, []byte(msg))

	if errSend != nil {
		fmt.Fprintf(s.errorW, "Unable to send email: %v\n", errSend)
		return errSend
	}

	return nil
}

// Levels return array of levels
func (s *SMTPConfig) Levels() []logrus.Level {
	return logrus.AllLevels
}

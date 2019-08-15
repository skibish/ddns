package notifier

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestTelegramInit(t *testing.T) {
	_, errConv := initTelegramNotifier("cfg")
	if errConv == nil {
		t.Error("Should be conversion error, but got nothing")
		return
	}

	m := make(map[interface{}]interface{})
	m["chat_id"] = 123
	_, errValue := initTelegramNotifier(m)
	if errValue == nil {
		t.Errorf("Should be error, because value is not a string")
		return
	}

	m = make(map[interface{}]interface{})
	m[123] = "123"
	_, errKey := initTelegramNotifier(m)
	if errKey == nil {
		t.Errorf("Should be error, because key is not a string")
		return
	}

	m = make(map[interface{}]interface{})
	m["chat_id"] = "123"
	tg, errUnexpected := initTelegramNotifier(m)
	if errUnexpected != nil {
		t.Error(errUnexpected)
		return
	}

	if tg.ChatID != "123" {
		t.Errorf("User should be %q, but got %q", "123", tg.ChatID)
		return
	}

}

func TestTelegramFire(t *testing.T) {

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	tg := TelegramConfig{
		Token:  "aaaa",
		ChatID: "1234",
		errorW: writer,
	}

	tg.send = func(msg string) error {
		return nil
	}

	entry := logrus.Entry{Message: "awesome message", Level: logrus.InfoLevel}
	errFire := tg.Fire(&entry)
	if errFire != nil {
		t.Errorf("handled an error on Telegram request: %q", errFire.Error())
		return
	}

	// test, got some error from Telegram
	tg.send = func(string) error {
		return errors.New("something went wrong")
	}

	errFire = tg.Fire(&entry)
	if errFire == nil {
		t.Error("Should be error, but it's OK")
		return
	}

	// test, DEBUG is ignored
	entry = logrus.Entry{Message: "ignored", Level: logrus.DebugLevel}
	shouldBeNil := tg.Fire(&entry)
	if shouldBeNil != nil {
		t.Errorf("Expected nil, on DEBUG, but got %v", shouldBeNil)
		return
	}
}

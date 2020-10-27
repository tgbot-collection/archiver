// archiver - handler
// 2020-10-27 21:45
// Benny <benny.think@gmail.com>

package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"net/url"
)

func startHandler(m *tb.Message) {

	_ = b.Notify(m.Chat, tb.Typing)
	_, _ = b.Send(m.Chat, "Hi! I'm the Internet Archive Wayback Machine bot.\n"+
		"You can send me any url and I'll save it for you.")

}

func aboutHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	_, _ = b.Send(m.Chat, "Wayback Machine bot by @BennyThink"+
		"GitHub: https://github.com/tgbot-collection/archiver")

}

func urlHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	_, err := url.ParseRequestURI(m.Text)
	if err != nil {
		log.Errorf("Bad url!%v", err)
		_, _ = b.Send(m.Chat, fmt.Sprintf("Your url %s seems to be invalid", m.Text))
	} else {
		result, err := taksSnapshot(m.Text)
		if err != nil {
			_, _ = b.Send(m.Chat, fmt.Sprintf("Failed! %v", err))
		} else {
			_, _ = b.Send(m.Chat, fmt.Sprintf("Success! %s", result))
		}

	}

}

func taksSnapshot(userUrl string) (string, error) {
	var body = url.Values{}
	body.Set("url", userUrl)
	body.Set("capture_all", "on")
	resp, err := http.PostForm(saveUrl, body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "success", nil
	} else {
		return fmt.Sprintf("fail %d", resp.StatusCode), errors.New(resp.Status)
	}

}

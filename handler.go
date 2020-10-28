// archiver - handler
// 2020-10-27 21:45
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/url"
)

func startHandler(m *tb.Message) {

	_ = b.Notify(m.Chat, tb.Typing)
	_, _ = b.Send(m.Chat, "Hi! I'm the Internet Archive Wayback Machine bot. https://archive.org/web/\n"+
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
		log.Errorf("Bad url! %v", err)
		_, _ = b.Send(m.Chat, fmt.Sprintf("Your url <pre>%s</pre> seems to be invalid", m.Text),
			&tb.SendOptions{ParseMode: tb.ModeHTML})
		return
	}
	replied, _ := b.Reply(m, "I've received your request. Please wait a second.")
	go takeSnapshot(m, replied)

}

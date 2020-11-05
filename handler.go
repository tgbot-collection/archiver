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

import "github.com/tgbot-collection/tgbot_ping"

func startHandler(m *tb.Message) {

	_ = b.Notify(m.Chat, tb.Typing)
	_, _ = b.Send(m.Chat, "Hi! I'm the [Internet Archive Wayback Machine bot](https://archive.org/web/).\n"+
		"You can send me any url and I'll archive it for you.", &tb.SendOptions{ParseMode: tb.ModeMarkdown})

}

func aboutHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	_, _ = b.Send(m.Chat, "Wayback Machine bot by @BennyThink"+
		"GitHub: https://github.com/tgbot-collection/archiver")

}

func pingHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	info := tgbot_ping.GetRuntime("botsrunner_archiver_1", "WaybackMachine Bot", "html")
	_, _ = b.Send(m.Chat, info, &tb.SendOptions{ParseMode: tb.ModeHTML})
}

func urlHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	_, err := url.ParseRequestURI(m.Text)
	if err != nil {
		log.Warnln("Invalid url.")
		_, _ = b.Send(m.Chat, fmt.Sprintf("Your url <pre>%s</pre> seems to be invalid", m.Text),
			&tb.SendOptions{ParseMode: tb.ModeHTML})
		return
	}
	replied, _ := b.Reply(m, "I've received your request. Please wait a second.")
	go takeSnapshot(m, replied)

}

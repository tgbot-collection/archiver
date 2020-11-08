// archiver - handler
// 2020-10-27 21:45
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"net/url"
	"time"
)

import (
	log "github.com/sirupsen/logrus"
	"github.com/tgbot-collection/tgbot_ping"
	tb "gopkg.in/tucnak/telebot.v2"
)

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
	// start to archive. Add more provider here ⬇️
	providers := []archiveProvider{&archiveOrg{}}
	for _, prov := range providers {
		go runner(m, replied, prov)
	}
}

func runner(m, replied *tb.Message, provider archiveProvider) {
	_, _ = b.Edit(replied, "Your archive request has been submitted.")

	_ = b.Notify(m.Chat, tb.UploadingDocument)
	html, err := provider.submit(m.Text)
	if err != nil {
		_, _ = b.Edit(replied, fmt.Sprintf("Request to %s failed:`%v`", provider, err),
			&tb.SendOptions{ParseMode: tb.ModeMarkdown})
		return
	}

	unique, err := provider.analysis(html)
	if err != nil {
		_, _ = b.Edit(replied, "Archive request has been submitted successfully. "+
			"But I'm unable to tell you current status. Generally this is okay to ignore.\nError: "+err.Error())
		return
	}

	_, _ = b.Edit(replied, "I'm trying to get the archive result for you...Please be patient.")
	var result string
	for i := 1; i <= 20; i++ {
		_ = b.Notify(m.Chat, tb.RecordingAudio)
		time.Sleep(time.Second * 7)
		result, err = provider.status(unique)
		// three-way handle
		if err != nil {
			log.Warnf("Refresh archive failed %v", err)
		} else if result != "" {
			_ = b.Notify(m.Chat, tb.Typing)
			// TODO if we're implementing more archive provider, we should consider reserve previous provider's result.
			_, _ = b.Edit(replied, result, &tb.SendOptions{ParseMode: tb.ModeMarkdown, DisableWebPagePreview: true})
			break
		} else if result == "" {
			msg := fmt.Sprintf("Refresh attempt %d/%d for %s", i, 20, m.Text)
			_, _ = b.Edit(replied, msg, &tb.SendOptions{DisableWebPagePreview: true})
		}
	}

	if result == "" {
		_, _ = b.Edit(replied, "Status operation timeout after 200s.")
	}
}

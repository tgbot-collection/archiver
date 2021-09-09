// archiver - handler
// 2020-10-27 21:45
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"regexp"
	"time"
)

import (
	log "github.com/sirupsen/logrus"
	"github.com/tgbot-collection/tgbot_ping"
	tb "gopkg.in/tucnak/telebot.v2"
)

func startHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	_, _ = b.Send(m.Chat, startText, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
}

func aboutHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	_, _ = b.Send(m.Chat, aboutText)
}

func pingHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	info := tgbot_ping.GetRuntime("botsrunner_archiver_1", "WaybackMachine Bot", "html")
	_, _ = b.Send(m.Chat, info, &tb.SendOptions{ParseMode: tb.ModeHTML})
}

func urlHandler(m *tb.Message) {
	_ = b.Notify(m.Chat, tb.Typing)
	replied, _ := b.Reply(m, Receive)
	providers := []archiveProvider{&archiveOrg{}}
	for _, prov := range providers {
		go runner(m, replied, prov)
	}
}

func runner(m, replied *tb.Message, provider archiveProvider) {
	re := regexp.MustCompile(`https?://.*`)
	urls := re.FindAllString(m.Text, -1)
	if len(urls) == 0 {
		_, _ = b.Edit(replied, InvalidRequest)
		return
	}

	for _, url := range urls {
		log.Infof("🗜️ Archiving %s", url)
		arc(m, replied, provider, url)
	}
	_, _ = b.Edit(replied, Finish)

}

func arc(m, replied *tb.Message, provider archiveProvider, url string) {
	_ = b.Notify(m.Chat, tb.UploadingDocument)
	html, err := provider.submit(url)
	if err != nil {
		_, _ = b.Edit(replied, fmt.Sprintf(ArchiveRequestFailed, provider, err))
		return
	}

	unique, err := provider.analysis(html)
	if err != nil {
		_, _ = b.Edit(replied, ArchiveNoResult+"\nError: "+err.Error())
		return
	}

	_, _ = b.Edit(replied, Processing)

	var result string
	for i := 1; i <= attempt; i++ {
		_ = b.Notify(m.Chat, tb.RecordingAudio)
		time.Sleep(sleep)
		result, err = provider.status(unique)
		// three-way handle
		if err != nil {
			log.Warnf("Refresh archive failed %v", err)
		} else if result != "" {
			_ = b.Notify(m.Chat, tb.Typing)
			_, _ = b.Send(m.Chat, result, &tb.SendOptions{ParseMode: tb.ModeMarkdown, DisableWebPagePreview: true})
			break
		} else if result == "" {
			msg := fmt.Sprintf(Updating, i, attempt, m.Text)
			_, _ = b.Edit(replied, msg, &tb.SendOptions{DisableWebPagePreview: true})
		}
	}

	if result == "" {
		_, _ = b.Edit(replied, ArchiveStatusTimeout)
	}
}

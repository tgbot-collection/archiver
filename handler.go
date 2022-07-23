// archiver - handler
// 2020-10-27 21:45
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

import (
	log "github.com/sirupsen/logrus"
	"github.com/tgbot-collection/tgbot_ping"
	tb "gopkg.in/telebot.v3"
)

func startHandler(c tb.Context) error {
	_ = b.Notify(c.Chat(), tb.Typing)
	_, _ = b.Send(c.Chat(), startText, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	return nil
}

func aboutHandler(c tb.Context) error {
	_ = b.Notify(c.Chat(), tb.Typing)
	_, _ = b.Send(c.Chat(), aboutText)
	return nil

}

func pingHandler(c tb.Context) error {
	_ = b.Notify(c.Chat(), tb.Typing)
	_ = b.Notify(c.Chat(), tb.Typing)
	info := tgbot_ping.GetRuntime("botsrunner_archiver_1", "WaybackMachine Bot", "html")
	ownerId, _ := strconv.ParseInt(os.Getenv("owner"), 10, 64)
	if c.Chat().ID == ownerId {
		info = fmt.Sprintf("%s\n Total URL archived %d", info, requestCount)
	}
	_, _ = b.Send(c.Chat(), info, &tb.SendOptions{ParseMode: tb.ModeHTML})
	return nil
}

func urlHandler(c tb.Context) error {
	_ = b.Notify(c.Chat(), tb.Typing)
	replied, _ := b.Reply(c.Message(), Receive)
	providers := []archiveProvider{&archiveOrg{}}
	for _, prov := range providers {
		go runner(c.Message(), replied, prov)
	}

	go takeScreenshot(c.Message().Text, c)
	return nil
}

func runner(m, replied *tb.Message, provider archiveProvider) {
	re := regexp.MustCompile(`https?://.*`)
	urls := re.FindAllString(m.Text, -1)
	if len(urls) == 0 {
		_, _ = b.Edit(replied, InvalidRequest)
		return
	}

	for _, url := range urls {
		requestCount += 1
		log.Infof("🗜️ Archiving %s", url)
		arc(m, replied, provider, url)
		time.Sleep(sleep / 2)
	}
	log.Infof(" %d jobs complted for %v", len(urls), m.Chat)
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
			log.Errorf("⚠️ %s refresh archive failed %v", url, err)
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
		log.Errorf(ArchiveStatusTimeout+" %s", url)
		_, _ = b.Edit(replied, ArchiveStatusTimeout)
	}
}

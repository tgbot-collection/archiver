// archiver - handler
// 2020-10-27 21:45
// Benny <benny.think@gmail.com>

package main

import (
	"errors"
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
	info := tgbot_ping.GetRuntime("botsrunner_archiver_1", "WaybackMachine Bot", "html")
	ownerId, _ := strconv.ParseInt(os.Getenv("owner"), 10, 64)
	if c.Chat().ID == ownerId {
		info = fmt.Sprintf("%s\n Total URL archived %d", info, requestCount)
	}
	_, _ = b.Send(c.Chat(), info, &tb.SendOptions{ParseMode: tb.ModeHTML})
	return nil
}

func mainEntrance(c tb.Context) error {
	_ = b.Notify(c.Chat(), tb.Typing)
	_ = b.Notify(c.Chat(), tb.Typing)
	user := getUser(c.Sender().ID)
	mode := user.Mode

	if mode == "ai" {
		if getChatsCount(c.Sender().ID) == 0 {
			// prepend a message
			addChat(c.Sender().ID, userRole, fmt.Sprintf("please take a look at this %s\n%s", user.Link, c.Message().Text))
		} else {
			addChat(c.Sender().ID, userRole, c.Message().Text)
		}
		aiResponse := askOpenAI(c.Sender().ID)
		addChat(c.Sender().ID, modelRole, aiResponse)
		return c.Send(aiResponse, tb.NoPreview)
	} else {
		return urlHandler(c)
	}

}

func stopAIHandler(c tb.Context) error {
	_ = b.Notify(c.Chat(), tb.Typing)
	disableAI(c.Sender().ID)
	rows := deleteChat(c.Sender().ID)
	return c.Send(fmt.Sprintf("AI mode disabled. %d chats deleted.", rows))
}

func buttonCallback(c tb.Context) error {
	var q = &tb.CallbackResponse{Text: "AI mode enabled."}
	enableAI(c.Sender().ID, c.Message().ReplyTo.Text)
	_ = c.Send("AI mode enabled. Please send your question. **Your chats will be saved in database until you use /stop to exit AI mode**",
		&tb.SendOptions{ParseMode: tb.ModeMarkdown})
	return c.Respond(q)
}

func urlHandler(c tb.Context) error {
	_ = b.Notify(c.Chat(), tb.Typing)
	replied, _ := b.Reply(c.Message(), Receive)
	providers := []archiveProvider{&archiveOrg{}}
	for _, prov := range providers {
		go archiveRunner(c.Message(), replied, prov)
	}

	go screenshotRunner(c)
	return nil
}

func screenshotRunner(c tb.Context) {
	urls, err := extractURL(c.Message().Text)
	if err != nil {
		return
	}
	for _, url := range urls {
		filename := takeScreenshot(url)
		_ = b.Notify(c.Chat(), tb.UploadingPhoto)
		p := &tb.Document{File: tb.FromDisk(filename), FileName: filename}
		_, _ = b.Reply(c.Message(), p)
		_ = os.Remove(filename)
	}
}

func archiveRunner(m, replied *tb.Message, provider archiveProvider) {
	urls, err := extractURL(m.Text)
	if err != nil {
		_, _ = b.Edit(replied, InvalidRequest)
		return
	}

	for _, url := range urls {
		requestCount += 1
		log.Infof("🗜️ Archiving %s", url)
		arc(m, replied, provider, url)
		time.Sleep(sleep / 2)
	}
	selector.Inline(selector.Row(btnPrev))
	_, _ = b.Edit(replied, Finish, selector)

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
		_, _ = b.Edit(replied, ArchiveStatusTimeout+" Please try here: https://web.archive.org/web/*/"+url)
	}
}

func extractURL(text string) ([]string, error) {
	re := regexp.MustCompile(`https?://.*`)
	urls := re.FindAllString(text, -1)
	if len(urls) == 0 {
		return []string{}, errors.New("No URL found")
	}
	return urls, nil
}

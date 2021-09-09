// archiver - config
// 2020-10-27 21:52
// Benny <benny.think@gmail.com>

package main

import (
	"os"
	"time"
)

var token = os.Getenv("TOKEN")

var attempt = 20
var sleep = time.Second * 5

var requestCount = 0

// post, param: url
var saveUrl = "https://web.archive.org/save/"

// get, https://web.archive.org/save/status/5d3157ab-6a03-4987-9847-b0e53ee84be9?_t=1603886202734
var statusUrl = "https://web.archive.org/save/status/"

type status struct {
	Status      string  `json:"status"`
	Timestamp   string  `json:"timestamp"`
	Duration    float32 `json:"duration_sec"`
	OriginalUrl string  `json:"original_url"`
}

const (
	Receive              = "🎬 Request received..."
	Processing           = "⌛️ Processing..."
	Updating             = "🗜️ Updating archive result...%d/%d for %s"
	Finish               = "🎉 Archive complete"
	InvalidRequest       = "❌ Your request was invalid"
	ArchiveNoResult      = "⚠️ Archive request has been submitted successfully. But I don't know result."
	ArchiveStatusTimeout = "❌ Archive status timeout"
	ArchiveRequestFailed = "❌ Request to %s failed: %v"
	aboutText            = "Wayback Machine bot by @BennyThink \nGitHub: https://github.com/tgbot-collection/archiver"
	startText            = "Hi! I'm the [Internet Archive Wayback Machine bot](https://archive.org/web/).\n" +
		"You can send me any url and I'll archive it for you."
)

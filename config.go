// archiver - config
// 2020-10-27 21:52
// Benny <benny.think@gmail.com>

package main

import "os"

var token = os.Getenv("TOKEN")

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

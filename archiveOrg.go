// archiver - archiveOrg
// 2020-11-08 16:21
// Benny <benny.think@gmail.com>

package main

// always keep standard library and 3rd separated
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

import (
	log "github.com/sirupsen/logrus"
)

// DO NOT involve any telegram bot objects here, such as  `b`, `*tb.Message`

type archiveOrg struct{}

func (a archiveOrg) submit(userUrl string) (html string, err error) {
	log.Infof("Requesting to archive.org %s", userUrl)

	var body = url.Values{}
	body.Set("url", userUrl)
	body.Set("capture_all", "on")

	req, err := http.NewRequest(http.MethodPost, saveUrl, bytes.NewBufferString(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("create archive request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Save Page Now returns a job_id as JSON when this header is present.
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Request to archive failed! %v", err)
		return "", fmt.Errorf("submit archive request: %w", err)
	}
	defer resp.Body.Close()

	tmp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read archive response: %w", err)
	}
	html = string(tmp)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("archive.org returned %s: %s", resp.Status, html)
	}

	log.Debugln("Requesting to archive.org has completed successfully.")
	return html, err
}

func (a archiveOrg) analysis(html string) (unique string, err error) {
	log.Debugln("Doing some analysis job....extracting unique UUID...")

	// The current Save Page Now API returns JSON, for example:
	// {"url":"https://example.com", "job_id":"..."}.
	var response struct {
		JobID string `json:"job_id"`
	}
	if json.Unmarshal([]byte(html), &response) == nil && response.JobID != "" {
		log.Debugln("Extraction success.")
		return response.JobID, nil
	}

	// Retain compatibility with the legacy HTML response.
	re := regexp.MustCompile(`spn\.watchJob\("(.+?)"`)
	result := re.FindStringSubmatch(html)

	if len(result) != 2 {
		err = fmt.Errorf("archive response contains no job_id")
		log.Errorf("Extract UUID failed! %v", err)
		return "", err
	}

	log.Debugln("Extraction success.")
	return result[1], nil
}

func (a archiveOrg) status(uuid string) (message string, err error) {
	reqUrl := fmt.Sprintf("%s%s?_t=%d", statusUrl, uuid, time.Now().Unix())
	resp, err := http.Get(reqUrl)
	if err != nil {
		return
	}

	var current status
	_ = json.NewDecoder(resp.Body).Decode(&current)

	if current.Status == "success" {
		message = fmt.Sprintf(
			`✅ %s, %.2fs. Click [here](%s) to visit your snapshot.`,
			current.Status,
			current.Duration,
			"https://web.archive.org/web/"+current.Timestamp+"/"+current.OriginalUrl,
		)
		log.Infof("✅ %s", current.OriginalUrl)
	} else {
		log.Infof("⌛️ %s - %s is  %s", current.OriginalUrl, time.Now(), current.Status)
	}

	_ = resp.Body.Close()
	return
}

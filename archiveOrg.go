// archiver - archiveOrg
// 2020-11-08 16:21
// Benny <benny.think@gmail.com>

package main

// always keep standard library and 3rd separated
import (
	"encoding/json"
	"errors"
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

	resp, err := http.PostForm(saveUrl, body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Errorf("Request to archive failed! %v", err)
		return "", err
	}

	tmp, _ := io.ReadAll(resp.Body)
	html = string(tmp)
	_ = resp.Body.Close()

	log.Debugln("Requesting to archive.org has completed successfully.")
	return html, err
}

func (a archiveOrg) analysis(html string) (unique string, err error) {
	log.Debugln("Doing some analysis job....extracting unique UUID...")

	uuid, err := __extractionUUID(html)
	if err != nil {
		log.Errorf("Extract UUID failed! %v", err)
		return "", err
	}

	log.Debugln("Extraction success.")
	return uuid, nil
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

func __extractionUUID(html string) (uuid string, err error) {
	re := regexp.MustCompile(`spn\.watchJob\("(.+?)"`)
	result := re.FindStringSubmatch(html)

	if len(result) != 2 {
		return "", errors.New(fmt.Sprintf("regex result is not equal to 2, %v", result))
	} else {
		return result[1], nil
	}

}

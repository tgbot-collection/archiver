package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func takeSnapshot(m, replied *tb.Message) {
	var userUrl = m.Text
	var body = url.Values{}
	body.Set("url", userUrl)
	body.Set("capture_all", "on")

	log.Infof("Requesting to %s with %v", saveUrl, userUrl)
	resp, err := http.PostForm(saveUrl, body)
	_, _ = b.Edit(replied, "Your archived request has been submitted.")

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Errorf("Request to archive failed! %v", err)
		_, _ = b.Edit(replied, fmt.Sprintf("Request failed: \n<pre>%v</pre>", err),
			&tb.SendOptions{ParseMode: tb.ModeHTML})
		return
	}

	// query the status
	html, _ := ioutil.ReadAll(resp.Body)
	uuid, err := extractionUuid(string(html))
	if err != nil {
		log.Errorf("Extract UUID failed! %v", err)
		_, _ = b.Edit(replied, "Archived request has been submitted successfully. "+
			"But I'm unable to tell you current status. Generally this is okay to ignore.\nError: "+err.Error())
		return
	}

	var snapResult = ""
	for i := 0; i <= 10; i++ {
		time.Sleep(time.Second * 10)
		st := retrieveStatus(uuid)
		if st != "" {
			snapResult = st
			break
		}
	}
	_ = resp.Body.Close()
	_, _ = b.Edit(replied, snapResult, &tb.SendOptions{ParseMode: tb.ModeHTML, DisableWebPagePreview: true})

}

func retrieveStatus(uuid string) (message string) {
	reqUrl := fmt.Sprintf("%s%s?_t=%d", statusUrl, uuid, time.Now().Unix())
	log.Infof("Getting new status from %s", reqUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		return
	}

	var current status
	_ = json.NewDecoder(resp.Body).Decode(&current)

	if current.Status == "success" {
		message = fmt.Sprintf(`%s, duration:%f
Click <a href="%s">here</a> to visit your snapshot.`, current.Status, current.Duration,
			"https://web.archive.org/web/"+current.Timestamp+"/"+current.OriginalUrl)
	} else {
		log.Infof("The result as of %s is still %s", time.Now(), current.Status)
	}

	_ = resp.Body.Close()
	return
}

func extractionUuid(html string) (uuid string, err error) {
	re := regexp.MustCompile(`spn\.watchJob\("(.+?)"`)
	result := re.FindStringSubmatch(html)

	if len(result) != 2 {
		return "", errors.New(fmt.Sprintf("regex result is not equal to 2, %v", result))
	} else {
		return result[1], nil
	}

}

// archiver - screenshot
// 2022-07-23 15:49
// Benny <benny.think@gmail.com>

package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	tb "gopkg.in/telebot.v3"
	"io/ioutil"
	"os"
)

var driver = os.Getenv("DRIVER")

var (
	chromeDriverPath = driver
	port             = 9515
)

func takeScreenshot(url string, c tb.Context) {
	// Start a WebDriver server instance
	opts := []selenium.ServiceOption{}
	selenium.SetDebug(false)

	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		log.Errorln(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	// headless
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--no-sandbox",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Errorln(err)
	}
	defer wd.Quit()

	if err := wd.Get(url); err != nil {
		log.Errorln(err)
	}

	//var width, _ = wd.ExecuteScript("return document.body.parentNode.scrollWidth", nil)
	var height, _ = wd.ExecuteScript("return document.body.parentNode.scrollHeight", nil)
	log.Infof("web page height: %.2f", height.(float64))
	err = wd.ResizeWindow("", 1920, int(height.(float64)))
	if err != nil {
		log.Errorln(err)
	}

	screenshot, _ := wd.Screenshot()
	// save screenshot to file
	//id := uuid.New()
	var filename = GetMD5Hash(url) + ".jpg"
	log.Infof("Saving screenshot to %s", filename)
	_ = ioutil.WriteFile(filename, screenshot, 0644)
	_ = b.Notify(c.Chat(), tb.UploadingPhoto)
	p := &tb.Document{File: tb.FromDisk(filename), FileName: filename}
	_, _ = b.Send(c.Chat(), p)
	log.Infof("Screenshot taken for %s", url)
	// delete file
	_ = os.Remove(filename)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

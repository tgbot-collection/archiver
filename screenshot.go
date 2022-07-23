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

var (
	chromeDriverPath = os.Getenv("DRIVER")
	port             = 9515
)

func takeScreenshot(url string, c tb.Context) {
	log.Infof("Taking screenshot for %s", url)
	// Start a WebDriver server instance
	var opts []selenium.ServiceOption
	selenium.SetDebug(false)

	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		log.Errorln(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
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

	var width, _ = wd.ExecuteScript("return document.body.parentNode.scrollWidth", nil)
	var height, _ = wd.ExecuteScript("return document.body.parentNode.scrollHeight", nil)
	var realHeight = int(height.(float64))
	var realWidth = int(width.(float64))

	log.Infof("web page width: %d, height: %d", realWidth, realHeight)
	err = wd.ResizeWindow("", 1920, realHeight)
	if err != nil {
		log.Errorln(err)
	}
	screenshot, _ := wd.Screenshot()
	log.Infof("screenshot size: %d", len(screenshot))
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

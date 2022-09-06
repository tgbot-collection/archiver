// archiver - screenshot
// 2022-07-23 15:49
// Benny <benny.think@gmail.com>

package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	tb "gopkg.in/telebot.v3"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"time"
)

const (
	label        = "Powered by https://t.me/wayback_machine_bot"
	browserWidth = 1440
	sleepTime    = 10 * time.Second
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

	caps := selenium.Capabilities{"browserName": "Chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
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
	beautify(url)
	_, _ = wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
	time.Sleep(5 * time.Second)

	var width, _ = wd.ExecuteScript("return document.body.parentNode.scrollWidth", nil)
	var height, _ = wd.ExecuteScript("return document.body.parentNode.scrollHeight", nil)
	var realHeight = int(height.(float64))
	var realWidth = int(width.(float64))

	log.Infof("web page width: %d, height: %d", realWidth, realHeight)

	// lazy loading
	const step = 1000
	var rounds = int(math.Ceil(float64(realHeight) / float64(step)))
	for i := 1; i <= rounds; i++ {
		_, _ = wd.ExecuteScript(fmt.Sprintf("window.scrollTo(0, %d)", i*step), nil)
		time.Sleep(300 * time.Millisecond)
	}

	_ = wd.ResizeWindow("", browserWidth, realHeight)
	_, _ = wd.ExecuteScript("window.scrollTo(0, 0)", nil)

	// wait for resource to load
	time.Sleep(sleepTime)

	screenshot, _ := wd.Screenshot()
	log.Infof("screenshot size: %d", len(screenshot))
	var filename = GetMD5Hash(url) + ".png"
	log.Infof("Saving screenshot to %s", filename)
	_ = os.WriteFile(filename, screenshot, 0644)
	addWatermark(filename)

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

func addWatermark(src string) {
	imgb, _ := os.Open(src)
	sourceImg, _ := png.Decode(imgb)
	waterMark := image.NewNRGBA(sourceImg.Bounds())

	fontList := []string{
		"/usr/share/fonts/TTF/SourceHanSans-VF.ttc",
		"/System/Library/Fonts/STHeiti Light.ttc",
		"/usr/share/fonts/TTF/OpenSans-Bold.ttf",
		"assets/Arial.ttf",
	}
	var font *truetype.Font
	for _, fontFile := range fontList {
		fontBytes, err := os.ReadFile(fontFile)
		if err == nil {
			log.Infof("Found font file %s", fontFile)
			font, _ = freetype.ParseFont(fontBytes)
			break
		}
	}

	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(font)
	f.SetFontSize(float64(browserWidth / 30))
	f.SetClip(sourceImg.Bounds())
	f.SetDst(waterMark)
	f.SetSrc(image.Black)

	width := waterMark.Bounds().Max.X
	height := 70
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			waterMark.Set(x, y, color.White)
		}
	}

	draw.Draw(waterMark, sourceImg.Bounds(), sourceImg, image.Pt(0, -70), draw.Src)
	pt := freetype.Pt(sourceImg.Bounds().Max.X/5, 50)
	_, _ = f.DrawString(label, pt)

	imgw, _ := os.Create(src)
	_ = png.Encode(imgw, waterMark)
	_ = imgb.Close()
	_ = imgw.Close()

}

//func main() {
//	takeScreenshot("https://www.baidu.com", nil)
//	addWatermark("1.png")
//}

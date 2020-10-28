// DailyGakki - gakki
// 2020-10-17 13:41
// Benny <benny.think@gmail.com>

package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

import (
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var b, err = tb.NewBot(tb.Settings{
	Token:  token,
	Poller: &tb.LongPoller{Timeout: 10 * time.Second},
})

func main() {
	if err != nil {
		log.Panicf("Please check your network or TOKEN! %v", err)
	}
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	Formatter := &log.TextFormatter{
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("[%s()]", f.Function), ""
		},
	}
	log.SetFormatter(Formatter)

	//  toilet  KeepMe.Run -f smblock
	banner := fmt.Sprintf(`
▌ ▌      ▌        ▌   ▙▗▌      ▌  ▗
▌▖▌▝▀▖▌ ▌▛▀▖▝▀▖▞▀▖▌▗▘ ▌▘▌▝▀▖▞▀▖▛▀▖▄ ▛▀▖▞▀▖
▙▚▌▞▀▌▚▄▌▌ ▌▞▀▌▌ ▖▛▚  ▌ ▌▞▀▌▌ ▖▌ ▌▐ ▌ ▌▛▀
▘ ▘▝▀▘▗▄▘▀▀ ▝▀▘▝▀ ▘ ▘ ▘ ▘▝▀▘▝▀ ▘ ▘▀▘▘ ▘▝▀▘
By %s 

Across the Great Wall, we can reach every corner in the world.
`, "BennyThink")

	fmt.Printf("\n %c[1;32m%s%c[0m\n\n", 0x1B, banner, 0x1B)

	b.Handle("/start", startHandler)
	b.Handle("/about", aboutHandler)
	b.Handle(tb.OnText, urlHandler)

	log.Infoln("I'm running...")
	b.Start()

}

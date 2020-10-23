//go:generate go get github.com/akavel/rsrc
//go:generate rsrc -ico mail.ico

package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/george-lewis/beeep"
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"github.com/sqweek/dialog"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type signal = struct{}

var (
	strLimit   string
	noMailIcon []byte
	mailIcon   []byte
	mailCh     chan int
	notifyCh   chan signal
	exitCh     [2]chan signal
	_config    config
)

func init() {
	var err error

	mailIcon, err = ioutil.ReadFile("mail.ico")
	if err != nil {
		dialog.Message("Couldn't load mail.ico").Title("Couldn't Load Icons").Error()
		panic(err)
	}

	noMailIcon, err = ioutil.ReadFile("nomail.ico")
	if err != nil {
		dialog.Message("Couldn't load nomail.ico").Title("Couldn't Load Icons").Error()
		panic(err)
	}

	color.Green("Images loaded")

	_config, err = readConfig(configFilename)

	if err != nil {
		dialog.Message("Couldn't load config file %s", configFilename).Title("Couldn't load config").Error()
		panic(err)
	}

	color.Green("Loaded config")

	mailCh = make(chan int)
	notifyCh = make(chan signal)

	for i := range exitCh {
		exitCh[i] = make(chan signal)
	}

	beeep.SetAppID("Reddit Inbox")

	strLimit = strconv.Itoa(_config.Limit)
}

func checkMail(bot reddit.Bot) (int, error) {
	h, err := bot.ListingWithParams("/message/unread", map[string]string{"limit": strLimit})

	if err != nil {
		color.Red("Reddit API Error: Could not get mail")
		return 0, err
	}

	return len(h.Messages), nil
}

func main() {
	color.Green("Startup")

	go systray.Run(onReady, onExit)

	bot, err := reddit.NewBotFromAgentFile("agent.txt", 0)
	if err != nil {
		dialog.Message("Did you put your agent.txt in the right location?").Title("Couldn't init Reddit API").Error()
		panic(err)
	}

	color.Green("Read agent file")

	go func() {
		_checkMail := func() {
			mail, err := checkMail(bot)
			if err == nil {
				mailCh <- mail
			}
		}
		sleepTime := time.Duration(_config.Interval) * time.Second
		for {
			select {
			case <-time.After(sleepTime):
				_checkMail()
			case <-notifyCh:
				_checkMail()
			case <-exitCh[0]:
				return
			}
		}
	}()

	cfg := graw.Config{CommentReplies: true, Messages: true, Mentions: true}

	handler := &mailer{bot: bot}

	stop, _, err := graw.Run(handler, bot, cfg)

	if err != nil {
		color.Red("Failed to start graw run: ", err)
	}

	color.Green("Started Reddit event listeners")

	<-exitCh[1]
	stop()

	// err = saveConfig(configFilename, _config)

	// if err != nil {
	// 	dialog.Message("Couldn't save config to %s", configFilename).Title("Couldn't save config").Error()
	// 	panic(err)
	// }

}

func onReady() {
	systray.SetIcon(noMailIcon)
	systray.SetTooltip("No Mail")

	inbox := systray.AddMenuItem("Inbox", "Go to your inbox")

	quit := systray.AddMenuItem("Quit", "Stop Reddit Mailer")

	color.Green("Systray ready")

	for {
		select {
		case <-inbox.ClickedCh:
			browser.OpenURL("https://www.reddit.com/message/unread")
		case <-quit.ClickedCh:
			systray.Quit()
		case c := <-mailCh:
			if c > 0 {
				systray.SetIcon(mailIcon)

				var plural string
				if c > 1 {
					plural = "s"
				} else {
					plural = ""
				}

				var plus string
				if c >= _config.Limit {
					plus = "+"
				} else {
					plus = ""
				}

				title := fmt.Sprintf("%d%s Message%s", c, plus, plural)
				systray.SetTitle(title)

				tooltip := fmt.Sprintf("You have %d%s message%s", c, plus, plural)
				systray.SetTooltip(tooltip)
			} else {
				systray.SetIcon(noMailIcon)
				systray.SetTitle("")
				systray.SetTooltip("No mail")
			}
		}
	}
}

func onExit() {
	for _, ch := range exitCh {
		ch <- signal{}
	}
}

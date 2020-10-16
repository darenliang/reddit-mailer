package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"github.com/turnage/graw/reddit"
)

const limit = 30

var (
	strLimit   string
	noMailIcon []byte
	mailIcon   []byte
	mailCh     chan int
	exitCh     chan bool
)

func init() {
	var err error

	mailIcon, err = readFile("mail.ico")
	if err != nil {
		panic(err)
	}

	noMailIcon, err = readFile("nomail.ico")
	if err != nil {
		panic(err)
	}

	color.Green("Images loaded")

	mailCh = make(chan int)
	exitCh = make(chan bool)

	strLimit = strconv.Itoa(limit)
}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func checkMail(bot reddit.Bot) (int, error) {
	h, err := bot.ListingWithParams("/message/unread", map[string]string{"limit": strLimit})
	if err != nil {
		return 0, err
	}

	return len(h.Messages), nil
}

func main() {
	color.Green("Startup")

	go systray.Run(onReady, onExit)

	bot, err := reddit.NewBotFromAgentFile("agent.txt", 0)
	if err != nil {
		panic(err)
	}

	color.Green("Reddit API initialized")

	mail, err := checkMail(bot)
	if err != nil {
		panic(err)
	}
	mailCh <- mail

	timer := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-timer.C:
			mail, err = checkMail(bot)
			if err != nil {
				panic(err)
			}
			mailCh <- mail
		case <-exitCh:
			timer.Stop()
			return
		}
	}
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
				if c >= limit {
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
				systray.SetTooltip("")
			}
		}
	}
}

func onExit() {
	exitCh <- true
}

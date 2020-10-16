package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/getlantern/systray"
	"github.com/turnage/graw/reddit"
)

var (
	noMailIcon []byte
	mailIcon   []byte
	mailCh     chan bool
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

	mailCh = make(chan bool)
	exitCh = make(chan bool)
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

func checkMail(bot reddit.Bot) (bool, error) {
	h, err := bot.ListingWithParams("/message/unread", map[string]string{"limit": "1"})
	if err != nil {
		return false, err
	}

	return len(h.Messages) > 0, nil
}

func main() {
	color.Green("Startup")

	go systray.Run(onReady, onExit)

	bot, err := reddit.NewBotFromAgentFile("agent.txt", 0)
	if err != nil {
		panic(err)
	}

	color.Green("Reddit API initialized")

	timer := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-timer.C:
			mail, err := checkMail(bot)
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
	systray.SetTitle("Reddit Mailer")

	quit := systray.AddMenuItem("Quit", "Stop Reddit Mailer")

	color.Green("Systray ready")

	for {
		select {
		case <-quit.ClickedCh:
			systray.Quit()
		case b := <-mailCh:
			if b {
				systray.SetIcon(mailIcon)
			} else {
				systray.SetIcon(noMailIcon)
			}
		}
	}
}

func onExit() {
	exitCh <- true
}

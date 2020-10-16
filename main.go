package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/getlantern/systray"
	"github.com/turnage/graw/reddit"
)

var nomail []byte
var mail []byte

var unreadCh chan int8
var exitCh chan int8

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

func checkMail(bot reddit.Bot) error {

	h, err := bot.ListingWithParams("/message/unread", map[string]string{"limit": "1"})

	if len(h.Messages) > 0 {
		unreadCh <- 1
	} else {
		unreadCh <- 0
	}

	if err != nil {
		return err
	}

	return nil

}

func main() {

	var err error

	color.Green("Startup")

	unreadCh = make(chan int8)
	exitCh = make(chan int8)

	mail, err = readFile("mail.ico")

	if err != nil {
		panic(err)
	}

	nomail, err = readFile("nomail.ico")

	if err != nil {
		panic(err)
	}

	color.Green("Images loaded")

	go systray.Run(onReady, onExit)

	var bot reddit.Bot

	bot, err = reddit.NewBotFromAgentFile("agent.txt", 0)

	if err != nil {
		panic(err)
	}

	color.Green("Reddit API initialized")

	checkMail(bot)

	timer := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-timer.C:
			err = checkMail(bot)
			if err != nil {
				panic(err)
			}
		case <-exitCh:
			return
		}
	}

}

func onReady() {

	systray.SetIcon(nomail)
	systray.SetTitle("Reddit Mailer")

	quit := systray.AddMenuItem("Quit", "Stop Reddit Mailer")

	color.Green("Systray ready")

	for {
		select {
		case <-quit.ClickedCh:
			systray.Quit()
		case c := <-unreadCh:
			if c == 1 {
				systray.SetIcon(mail)
			} else {
				systray.SetIcon(nomail)
			}
		}
	}

}

func onExit() {
	exitCh <- 1
}

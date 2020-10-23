package main

import (
	"fmt"
	"github.com/george-lewis/beeep"
	"github.com/turnage/graw/reddit"
)

type mailer struct {
	bot reddit.Bot
}

func (m *mailer) CommentReply(reply *reddit.Message) error {
	title := fmt.Sprintf("/u/%s replied to you", reply.Author)
	return processEvent(title, reply.Body, _config.CommentReplies)
}

func (m *mailer) Message(msg *reddit.Message) error {
	title := fmt.Sprintf("/u/%s sent you a message", msg.Author)
	return processEvent(title, msg.Body, _config.Messages)
}

func (m *mailer) PostReply(reply *reddit.Message) error {
	title := fmt.Sprintf("/u/%s replied to your post", reply.Author)
	return processEvent(title, reply.Body, _config.PostReplies)
}

func (m *mailer) Mention(mention *reddit.Message) error {
	title := fmt.Sprintf("/u/%s mentioned you", mention.Author)
	return processEvent(title, mention.Body, _config.Mentions)
}

func processEvent(title string, body string, notify bool) error {
	notifyCh <- signal{}
	if _config.Notifications && notify {
		beeep.Notify(title, body, "mail.ico")
	}
	return nil
}

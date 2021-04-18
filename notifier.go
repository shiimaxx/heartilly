package main

import "github.com/slack-go/slack"

type Notifier interface {
	Notify(Message) error
}

type SlackNotifier struct {
	Channel string
	Client  *slack.Client
}

func NewSlackNotifier(token, channel string) *SlackNotifier {
	return &SlackNotifier{
		Channel: channel,
		Client:  slack.New(token),
	}
}

func (s *SlackNotifier) Notify(msg Message) error {
	_, _, err := s.Client.PostMessage(
		s.Channel,
		slack.MsgOptionAttachments(
			slack.Attachment{
				Color: s.color(msg.StatusType),
				Text:  msg.Text,
			},
		),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SlackNotifier) color(status Status) string {
	switch {
	case status == OK:
		return "good"
	case status == Alert:
		return "danger"
	default:
		return "#808080"
	}
}

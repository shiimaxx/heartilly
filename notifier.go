package main

import "github.com/slack-go/slack"

type Notifier interface {
	Notify(string) error
}

type SlackNotifier struct {
	Channel string
	Client  *slack.Client
}

func NewSlackNotifier(token, channel string) *SlackNotifier {
	return &SlackNotifier{
		Channel: channel,
		Client: slack.New(token),
	}
}

func (s *SlackNotifier) Notify(msg string) error {
	_, _, err := s.Client.PostMessage(
		s.Channel,
		slack.MsgOptionText(msg, false),
	)
	if err != nil {
		return err
	}
	return nil
}

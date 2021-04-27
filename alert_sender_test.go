package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type NotifierMock struct {
	mock.Mock
}

func (m *NotifierMock) Notify(msg Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func TestSetNotifier(t *testing.T) {
	alertSender := &AlertSender{}
	alertSender.SetNotifier(NewSlackNotifier("token", "channel"))
	alertSender.SetNotifier(NewSlackNotifier("token", "channel"))

	want := 2
	got := len(alertSender.Notifiers)
	assert.Equal(t, want, got)
}

func TestRun(t *testing.T) {
	messageCh := make(chan Message)
	errCh := make(chan error)

	dummyMessage := Message{"dummy message", Critical}
	notifierMock := new(NotifierMock)

	notifierMock.On("Notify", dummyMessage).Return(nil)

	alertSender := &AlertSender{
		Notifiers: []Notifier{notifierMock},
		MessageCh: messageCh,
		ErrCh:     errCh,
	}
	go alertSender.Run()

	messageCh <- dummyMessage
	time.Sleep(1 * time.Millisecond)

	notifierMock.AssertExpectations(t)
}

func TestRun_notify_error(t *testing.T) {
	messageCh := make(chan Message)
	errCh := make(chan error)

	dummyMessage := Message{"dummy message", Critical}
	notifierMock := new(NotifierMock)

	notifierMock.On("Notify", dummyMessage).Return(fmt.Errorf("error"))

	alertSender := &AlertSender{
		Notifiers: []Notifier{notifierMock},
		MessageCh: messageCh,
		ErrCh:     errCh,
	}
	go alertSender.Run()

	messageCh <- dummyMessage
	time.Sleep(1 * time.Millisecond)

	notifierMock.AssertExpectations(t)

	err := <-errCh
	assert.NotNil(t, err)
}

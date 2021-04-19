package main

type AlertSender struct {
	Notifiers []Notifier

	MessageCh <-chan Message
	ErrCh     chan<- error
}

func (as *AlertSender) SetNotifier(n Notifier) {
	as.Notifiers = append(as.Notifiers, n)
}

func (as *AlertSender) Run() {
	for {
		msg := <-as.MessageCh

		for _, notifier := range as.Notifiers {
			if err := notifier.Notify(msg); err != nil {
				as.ErrCh <- err
			}
		}
	}
}

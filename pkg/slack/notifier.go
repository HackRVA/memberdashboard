package slack

type Notifier struct{}

func (s Notifier) Send(msg string) {
	Send(msg)
}

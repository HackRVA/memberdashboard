package slack

type Notifier struct {
	WebHookURL string
}

func (s Notifier) Send(msg string) {
	Send(s.WebHookURL, msg)
}

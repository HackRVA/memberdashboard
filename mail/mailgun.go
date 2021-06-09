package mail

import (
	"context"
	"net/smtp"
	"time"

	"memberserver/config"

	"github.com/jordan-wright/email"
	"github.com/mailgun/mailgun-go/v4"
	log "github.com/sirupsen/logrus"
)

// Provider has config information for connecting to mailgun
type Provider struct {
	from     string
	user     string
	password string
	URL      string
	Key      string
}

// Setup attaches config information to mailProvider object
func Setup() (*Provider, error) {
	mp := &Provider{}
	c, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	mp.URL = c.MailgunURL
	mp.Key = c.MailgunKey
	mp.from = c.MailgunFromAddress
	mp.user = c.MailgunUser
	mp.password = c.MailgunPassword

	return mp, nil
}

// SendSMTP - send a simple email via smtp
func (mp Provider) SendSMTP(address, subject, text string) (string, error) {
	e := email.NewEmail()
	e.From = "noreply <" + mp.from + ">"
	e.To = []string{address}
	e.Subject = subject
	e.Text = []byte(text)
	err := e.Send("smtp.mailgun.org:587", smtp.PlainAuth("", mp.user, mp.password, "smtp.mailgun.org"))
	return "_", err
}

func (mp Provider) SendPlainTextMail(address, subject, text string) (string, error) {
	mg := mailgun.NewMailgun(mp.URL, mp.Key)
	m := mg.NewMessage(
		"noreply <"+mp.from+">",
		subject,
		text,
		address,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func (mp Provider) SendHtmlMail(address, subject, html string) (string, error) {
	mg := mailgun.NewMailgun(mp.URL, mp.Key)
	m := mg.NewMessage(
		"noreply <"+mp.from+">",
		subject,
		"Testing some Mailgun awesomeness!",
		address,
	)
	// m.AddCC("info@hackrva.org")
	// m.AddBCC("bar@example.com")
	m.SetHtml(html)
	// m.AddAttachment("files/test.jpg")
	// m.AddAttachment("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

package mail

import (
	"context"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/jordan-wright/email"
	"github.com/mailgun/mailgun-go/v4"

	"github.com/dfirebaugh/memberserver/config"
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
	println(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))
	c, err := config.Load(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))

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

// SendSimpleMessage - sends an email via the api
//  I haven't been able to get this to work due to credentials issues.
//  I'm assuming I'm doing something wrong in the control panel
//  eventually this will be the better way to send email
func (mp Provider) SendSimpleMessage(address, subject, text string) (string, error) {
	mg := mailgun.NewMailgun(mp.URL, mp.Key)
	m := mg.NewMessage(
		"noreply <"+mp.from+">",
		subject,
		text,
		address,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	println("attempting to send simple mail")

	_, id, err := mg.Send(ctx, m)
	return id, err
}

// SendComplexMessage - sends an email via the api but allows for html body to be attached
//  I haven't been able to get this to work due to credentials issues.
//  I'm assuming I'm doing something wrong in the control panel
//  eventually this will be the better way to send email
func (mp Provider) SendComplexMessage(address, subject, html string) (string, error) {
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

	println("attempting to send complex mail")

	_, id, err := mg.Send(ctx, m)
	return id, err
}

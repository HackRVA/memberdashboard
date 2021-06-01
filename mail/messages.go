package mail

import (
	"bytes"
	"html/template"
	"memberserver/config"

	log "github.com/sirupsen/logrus"
)

func SendGracePeriodMessageToLeadership(recipient string, member interface{}) {
	infoAddress := "info@hackrva.org"
	SendTemplatedEmail("pending_revokation_leadership.html.tmpl", infoAddress, "hackRVA Grace Period", member)
}

func SendGracePeriodMessage(recipient string, member interface{}) {
	SendTemplatedEmail("pending_revokation_member.html.tmpl", recipient, "hackRVA Grace Period", member)
}

func SendRevokedEmail(recipient string, member interface{}) {
	SendTemplatedEmail("access_revoked.html.tmpl", recipient, "hackRVA Grace Period", member)
}

func SendRevokedEmailToLeadership(recipient string, member interface{}) {
	SendTemplatedEmail("access_revoked_leadership.html.tmpl", recipient, "hackRVA Grace Period", member)
}

func SendIPHasChanged(newIpAddress string) {
	recipient := "info@hackrva.org"
	model := struct {
		IpAddress string
	}{
		IpAddress: newIpAddress}
	SendTemplatedEmail("ip_changed.html.tmpl", recipient, "IP Address Changed", model)
}

func generateEmailContent(templatePath string, model interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Errorf("Error loading template %v", err)
		return "", err
	}
	tmpl.Option("missingkey=error")
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, model)
	if err != nil {
		log.Errorf("Error generating content %v", err)
		return "", err
	}
	return tpl.String(), nil
}

func SendTemplatedEmail(templateName string, to string, subject string, model interface{}) {
	conf, _ := config.Load()

	if !conf.EnableInfoEmails {
		log.Info("email not enabled")
		return
	}

	mp, err := Setup()
	if err != nil {
		log.Errorf("error setting up mailprovider when attempting to send email notification")
		return
	}

	if len(conf.EmailOverrideAddress) > 0 {
		to = conf.EmailOverrideAddress
	}

	content, err := generateEmailContent("./templates/"+templateName, model)
	if err != nil {
		log.Errorf("Error generating email contnent. Error: %v", err)
		return
	}

	_, err = mp.SendComplexMessage(to, subject, content)
	if err != nil {
		log.Errorf("Error sending mail %v", err)
	}
}

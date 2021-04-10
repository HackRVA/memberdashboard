package mail

import (
	"memberserver/config"

	log "github.com/sirupsen/logrus"
)

func SendGracePeriodMessageToLeadership(address string) {
	infoAddress := "info@hackrva.org"

	conf, _ := config.Load()
	if !conf.EnableInfoEmails {
		return
	}

	mp, err := Setup()
	if err != nil {
		log.Errorf("error setting up mailprovider when attempting to send email notification")
	}

	if len(conf.EmailOverrideAddress) > 0 {
		infoAddress = conf.EmailOverrideAddress
	}

	mp.SendSMTP(infoAddress, address+": hackrva grace period", address+" membership is in a grace period.  \n\nIf a payment isn't received, their membership will be revoked.")
}

func SendGracePeriodMessage(address string) {
	conf, _ := config.Load()
	if !conf.EnableNotificationEmailsToMembers {
		return
	}
	mp, err := Setup()
	if err != nil {
		log.Errorf("error setting up mailprovider when attempting to send email notification")
	}

	if len(conf.EmailOverrideAddress) > 0 {
		address = conf.EmailOverrideAddress
	}

	mp.SendSMTP(address, "hackrva grace period", "This is an automated message.\n\n You're membership is in a grace period.  Please try to pay your hackrva membership dues soon.  If you have concerns, please reach out to info@hackrva.org")
}

func SendRevokedEmail(address string) {
	conf, _ := config.Load()
	if !conf.EnableNotificationEmailsToMembers {
		return
	}
	mp, err := Setup()
	if err != nil {
		log.Errorf("error setting up mailprovider when attempting to send email notification")
	}

	if len(conf.EmailOverrideAddress) > 0 {
		address = conf.EmailOverrideAddress
	}

	mp.SendSMTP(address, "hackrva membership revoked", "Unfortunately, hackrva hasn't received your membership dues.  Your membership has been revoked until a payment is received.  Please reach out to us if you have any concerns.")
}

func SendIPHasChanged() {
	conf, _ := config.Load()
	if !conf.EnableInfoEmails {
		return
	}

	mp, err := Setup()
	if err != nil {
		log.Errorf("error setting up mailprovider when attempting to send email notification")
	}

	infoAddress := "info@hackrva.org"

	if len(conf.EmailOverrideAddress) > 0 {
		infoAddress = conf.EmailOverrideAddress
	}

	mp.SendSMTP(infoAddress, "hackrva's ip address has changed", "HackRVAs IP address has changed.  This is significant because the database for the member dashboard is IP whitelisted.  For the dashboard to work, someone will need to update the whitelisting.")
}

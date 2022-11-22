package mail

import (
	"errors"
	"time"

	"github.com/HackRVA/memberserver/internal/services/config"

	"github.com/HackRVA/memberserver/internal/models"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type CommunicationTemplate string

const (
	AccessRevokedMember         CommunicationTemplate = "AccessRevokedMember"
	AccessRevokedLeadership     CommunicationTemplate = "AccessRevokedLeadership"
	IpChanged                   CommunicationTemplate = "IpChanged"
	PendingRevokationLeadership CommunicationTemplate = "PendingRevokationLeadership"
	PendingRevokationMember     CommunicationTemplate = "PendingRevokationMember"
	Welcome                     CommunicationTemplate = "Welcome"
)

// String converts CommunicationTemplate to a string
func (c CommunicationTemplate) String() string {
	return string(c)
}

type mailer struct {
	db        CommunicationDal
	m         MailApi
	config    config.Config
	generator templateGenerator
}

type MailApi interface {
	SendHtmlMail(address, subject, body string) (string, error)
	SendPlainTextMail(address, subject, content string) (string, error)
}

type CommunicationDal interface {
	GetMemberByEmail(memberEmail string) (models.Member, error)
	GetCommunication(communication string) (models.Communication, error)
	LogCommunication(communicationId int, memberId string) error
	GetMostRecentCommunicationToMember(memberId string, commId int) (time.Time, error)
}

func NewMailer(db CommunicationDal, m MailApi, config config.Config) *mailer {
	mailer := mailer{
		db,
		m,
		config,
		fileTemplateGenerator{}}
	return &mailer
}

func (m *mailer) SendCommunication(communication CommunicationTemplate, recipient string, model interface{}) (bool, error) {
	memberExists := true
	member, err := m.db.GetMemberByEmail(recipient)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			memberExists = false
		} else {
			return false, err
		}
	}

	if memberExists && !m.config.EnableNotificationEmailsToMembers {
		log.Printf("Communication %v not sent to %v because notifications to members is not enabled", communication.String(), recipient)
		return false, nil
	} else if !memberExists && !m.config.EnableInfoEmails {
		log.Printf("Communication %v not sent to %v because info emails not enabled", communication.String(), recipient)
		return false, nil
	}

	c, err := m.db.GetCommunication(communication.String())
	if err != nil {
		log.Printf("%v not found. Err: %v", communication.String(), err)
		return false, err
	}

	if memberExists && m.IsThrottled(c, member) {
		log.Printf("Communication %v not sent to %v due to throttling", communication.String(), recipient)
		return false, nil
	}

	content, err := m.generator.generateEmailContent("./templates/"+c.Template, model)
	if err != nil {
		log.Errorf("Error generating email content. Error: %v", err)
		return false, err
	}

	if len(m.config.EmailOverrideAddress) > 0 {
		recipient = m.config.EmailOverrideAddress
	}

	_, err = m.m.SendHtmlMail(recipient, c.Subject, content)
	if err != nil {
		log.Printf("Failed to send mail to %v.  Err: %v", recipient, err)
		return false, err
	}

	if memberExists {
		m.db.LogCommunication(c.ID, member.ID)
	}

	return true, nil
}

func (m *mailer) IsThrottled(c models.Communication, member models.Member) bool {

	if c.FrequencyThrottle > 0 {
		last, err := m.db.GetMostRecentCommunicationToMember(member.ID, c.ID)
		if err != nil {
			return false
		}
		difference := time.Since(last).Hours() / 24
		if difference < float64(c.FrequencyThrottle) {
			return true
		}
	}
	return false
}

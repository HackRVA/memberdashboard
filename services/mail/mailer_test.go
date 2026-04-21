package mail

import (
	"context"
	"testing"
	"time"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/models"

	"github.com/jackc/pgx/v4"
)

func TestSendMessageToNonMemberWithoutLogging(t *testing.T) {
	ctx := context.Background()
	db := dbMock{}
	m := mailApiMock{}
	c, _ := config.Load()
	c.EnableInfoEmails = true

	mailer := NewMailer(&db, &m, c)
	mailer.generator = generatorMock{}

	db.memberError = pgx.ErrNoRows

	sent, err := mailer.SendCommunication(ctx, AccessRevokedLeadership, c.AdminEmail, memberModel)
	if err != nil {
		t.Errorf("Error sending communication %v", err)
	}
	if !sent {
		t.Error("Mail not sent")
	}

	if db.logCommunicationCalled {
		t.Error("Communications to non-members should not be logged")
	}
}

func TestSendMessageToMemberShouldLog(t *testing.T) {
	ctx := context.Background()
	db := dbMock{}
	m := mailApiMock{}
	c, _ := config.Load()
	c.EnableNotificationEmailsToMembers = true

	mailer := NewMailer(&db, &m, c)
	mailer.generator = generatorMock{}

	sent, err := mailer.SendCommunication(ctx, AccessRevokedLeadership, "member@hackrva.org", memberModel)
	if err != nil {
		t.Errorf("Error sending communication %v", err)
	}
	if !sent {
		t.Error("Mail not sent")
	}

	if !db.logCommunicationCalled {
		t.Error("Communications to members should be logged")
	}
}

func TestSendMessageToShouldThrottle(t *testing.T) {
	ctx := context.Background()
	db := dbMock{}
	c, _ := config.Load()
	c.EnableNotificationEmailsToMembers = true
	m := mailApiMock{}

	mailer := NewMailer(&db, &m, c)
	mailer.generator = generatorMock{}
	db.communicationResult.FrequencyThrottle = 10
	db.mostRecentCommResult = time.Now().AddDate(0, 0, -5)

	sent, err := mailer.SendCommunication(ctx, AccessRevokedLeadership, c.AdminEmail, memberModel)
	if err != nil {
		t.Errorf("Error sending communication %v", err)
	}
	if sent {
		t.Error("Mail should not be sent since it was within throttle")
	}

	if db.logCommunicationCalled {
		t.Error("Log should not be created since communication was not sent")
	}
}

func TestEnableMemberEmailsSetShouldSendMemberEmails(t *testing.T) {
	ctx := context.Background()
	db := dbMock{}
	m := mailApiMock{}
	c, _ := config.Load()
	c.EnableNotificationEmailsToMembers = true

	mailer := NewMailer(&db, &m, c)
	mailer.generator = generatorMock{}

	sent, err := mailer.SendCommunication(ctx, AccessRevokedMember, "member@email.com", memberModel)
	if err != nil {
		t.Errorf("Error sending communication %v", err)
	}
	if !sent {
		t.Error("Mail not sent")
	}
	if !m.MailSent {
		t.Error("Mail not sent")
	}
}

func TestEnableMemberEmailsUnsetShouldNotSendMemberEmails(t *testing.T) {
	ctx := context.Background()
	db := dbMock{}
	m := mailApiMock{}
	c, _ := config.Load()
	c.EnableNotificationEmailsToMembers = false

	mailer := NewMailer(&db, &m, c)
	mailer.generator = generatorMock{}

	sent, err := mailer.SendCommunication(ctx, AccessRevokedMember, "member@email.com", memberModel)
	if err != nil {
		t.Errorf("Error sending communication %v", err)
	}
	if sent {
		t.Error("Mail should not be sent")
	}
	if m.MailSent {
		t.Error("Mail should not be sent")
	}
}

func TestEnableInfoEmailsSetShouldSendInfoEmails(t *testing.T) {
	ctx := context.Background()
	db := dbMock{}
	m := mailApiMock{}
	c, _ := config.Load()
	c.EnableInfoEmails = true
	db.memberError = pgx.ErrNoRows

	mailer := NewMailer(&db, &m, c)
	mailer.generator = generatorMock{}

	sent, err := mailer.SendCommunication(ctx, AccessRevokedMember, c.AdminEmail, memberModel)
	if err != nil {
		t.Errorf("Error sending communication %v", err)
	}
	if !sent {
		t.Error("Mail should not be sent")
	}
	if !m.MailSent {
		t.Error("Mail should not be sent")
	}
}

func TestEnableInfoEmailsUnsetShouldNotSendInfoEmails(t *testing.T) {
	ctx := context.Background()
	db := dbMock{}
	m := mailApiMock{}
	c, _ := config.Load()
	c.EnableInfoEmails = false
	db.memberError = pgx.ErrNoRows

	mailer := NewMailer(&db, &m, c)
	mailer.generator = generatorMock{}

	sent, err := mailer.SendCommunication(ctx, AccessRevokedMember, c.AdminEmail, memberModel)
	if err != nil {
		t.Errorf("Error sending communication %v", err)
	}
	if sent {
		t.Error("Mail should not be sent")
	}
	if m.MailSent {
		t.Error("Mail should not be sent")
	}
}

type dbMock struct {
	memberResult           models.Member
	memberError            error
	communicationResult    models.Communication
	communicatonError      error
	mostRecentCommResult   time.Time
	mostRecentCommError    error
	logCommunicationCalled bool
}

func (m *dbMock) GetMemberByEmail(ctx context.Context, memberEmail string) (models.Member, error) {
	return m.memberResult, m.memberError
}

func (m *dbMock) GetCommunication(ctx context.Context, communication string) (models.Communication, error) {
	return m.communicationResult, m.communicatonError
}

func (m *dbMock) LogCommunication(ctx context.Context, communicationId int, memberId string) error {
	m.logCommunicationCalled = true
	return nil
}

func (m *dbMock) GetMostRecentCommunicationToMember(ctx context.Context, memberId string, commId int) (time.Time, error) {
	return m.mostRecentCommResult, m.mostRecentCommError
}

type mailApiMock struct {
	MailSent bool
}

func (m *mailApiMock) SendHtmlMail(address, subject, body string) (string, error) {
	m.MailSent = true
	return "", nil
}

func (m *mailApiMock) SendPlainTextMail(address, subject, content string) (string, error) {
	m.MailSent = true
	return "", nil
}

type generatorMock struct{}

func (generatorMock) generateEmailContent(templateSource string, model interface{}) (string, error) {
	return "", nil
}

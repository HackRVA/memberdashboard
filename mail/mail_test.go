package mail

import (
	"testing"
)

func TestIpChangedTemplate(t *testing.T) {
	model := struct {
		IpAddress string
	}{
		IpAddress: "127.0.0.1",
	}
	content, err := generateEmailContent("../templates/ip_changed.html.tmpl", model)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestAccessRevokedLeadershipTemplate(t *testing.T) {
	model := struct {
		Name  string
		Email string
	}{
		Name:  "Member Name",
		Email: "member@email.com",
	}
	content, err := generateEmailContent("../templates/access_revoked_leadership.html.tmpl", model)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestAccessRevokedTemplate(t *testing.T) {
	model := struct {
		Name  string
		Email string
	}{
		Name:  "Member Name",
		Email: "member@email.com",
	}
	content, err := generateEmailContent("../templates/access_revoked.html.tmpl", model)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestPendingRevokationLeadershipTemplate(t *testing.T) {
	model := struct {
		Name  string
		Email string
	}{
		Name:  "Member Name",
		Email: "member@email.com",
	}
	content, err := generateEmailContent("../templates/pending_revokation_leadership.html.tmpl", model)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestPendingRevokationMemberTemplate(t *testing.T) {
	model := struct {
		Name  string
		Email string
	}{
		Name:  "Member Name",
		Email: "member@email.com",
	}
	content, err := generateEmailContent("../templates/pending_revokation_member.html.tmpl", model)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestWelcomeTemplate(t *testing.T) {
	model := struct {
		Name  string
		Email string
	}{
		Name:  "Member Name",
		Email: "member@email.com",
	}
	content, err := generateEmailContent("../templates/welcome.html.tmpl", model)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

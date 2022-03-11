package mail

import (
	"testing"
)

var generator fileTemplateGenerator = fileTemplateGenerator{}
var memberModel = struct {
	Name  string
	Email string
}{
	Name:  "Member Name",
	Email: "member@email.com",
}

func TestIpChangedTemplate(t *testing.T) {
	ipModel := struct {
		IpAddress string
	}{
		IpAddress: "127.0.0.1",
	}
	content, err := generator.generateEmailContent("../../../templates/ip_changed.html.tmpl", ipModel)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestAccessRevokedLeadershipTemplate(t *testing.T) {
	content, err := generator.generateEmailContent("../../../templates/access_revoked_leadership.html.tmpl", memberModel)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestAccessRevokedTemplate(t *testing.T) {
	content, err := generator.generateEmailContent("../../../templates/access_revoked.html.tmpl", memberModel)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestPendingRevokationLeadershipTemplate(t *testing.T) {
	content, err := generator.generateEmailContent("../../../templates/pending_revokation_leadership.html.tmpl", memberModel)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestPendingRevokationMemberTemplate(t *testing.T) {
	content, err := generator.generateEmailContent("../../../templates/pending_revokation_member.html.tmpl", memberModel)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

func TestWelcomeTemplate(t *testing.T) {
	content, err := generator.generateEmailContent("../../../templates/welcome.html.tmpl", memberModel)
	if err != nil {
		t.Fatalf("Failed to generate content. %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("Failed to generate content.  Result is empty")
	}
}

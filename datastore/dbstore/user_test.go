package dbstore

import (
	"testing"

	"github.com/HackRVA/memberserver/models"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

func TestGetUser_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT email from membership\.users where email=\$1`).
		WithArgs("alice@example.com").
		WillReturnRows(pgxmock.NewRows([]string{"email"}).AddRow("alice@example.com"))

	got, err := db.GetUser("Alice@Example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Email != "alice@example.com" {
		t.Errorf("Email = %q, want alice@example.com", got.Email)
	}
}

func TestUserSignin_UnauthorizedOnNoRows(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT password from membership\.users where email=\$1`).
		WithArgs("alice@example.com").
		WillReturnError(pgx.ErrNoRows)

	err := db.UserSignin("Alice@Example.com", "password")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRegisterUser_EmptyPassword(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	err := db.RegisterUser(models.Credentials{Email: "alice@example.com", Password: ""})
	if err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestRegisterUser_EmptyEmail(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	err := db.RegisterUser(models.Credentials{Email: "", Password: "secret"})
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestRegisterUser_MemberLookupFails(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs("ghost@example.com").
		WillReturnError(pgx.ErrNoRows)

	err := db.RegisterUser(models.Credentials{Email: "ghost@example.com", Password: "secret"})
	if err == nil {
		t.Fatal("expected error -- member must exist before register")
	}
}

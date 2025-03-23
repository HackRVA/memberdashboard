package dbstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/HackRVA/memberserver/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser register a user in the db
func (db *DatabaseStore) RegisterUser(creds models.Credentials) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	if len(creds.Password) == 0 {
		return fmt.Errorf("not a valid password")
	}

	if len(creds.Email) == 0 {
		return fmt.Errorf("not a valid email")
	}

	// require the user to be a member
	_, err = db.GetMemberByEmail(creds.Email)
	if err != nil {
		return err
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		return err
	}

	// Next, insert the email, along with the hashed password into the database
	rows, err := dbPool.Query(context.Background(), userDbMethod.registerUser(), creds.Email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("registerUser failed: %s", err)
	}

	defer rows.Close()

	return nil
}

// UserSignin - user login
func (db *DatabaseStore) UserSignin(email string, password string) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	// We create another instance of `Credentials` to store the credentials we get from the database
	storedCreds := &models.Credentials{}

	// Get the existing entry present in the database for the given user
	row := dbPool.QueryRow(context.Background(), userDbMethod.getUserPassword(), strings.ToLower(email)).Scan(&storedCreds.Password)
	if row == pgx.ErrNoRows {
		return fmt.Errorf("Unauthorized")
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(password)); err != nil {
		// If the two passwords don't match, return a 401 status
		return fmt.Errorf("unauthorized: %s", err)
	}

	return nil
}

// GetUser returns the currently logged in user
func (db *DatabaseStore) GetUser(email string) (models.UserResponse, error) {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	var userResponse models.UserResponse

	row := dbPool.QueryRow(context.Background(), userDbMethod.getUser(), strings.ToLower(email)).Scan(&userResponse.Email)
	if row == pgx.ErrNoRows {
		return userResponse, fmt.Errorf("error getting user")
	}
	return userResponse, nil
}

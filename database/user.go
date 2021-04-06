package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// Credentials Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	// Password - the user's password
	// required: true
	// example: string
	Password string `json:"password"`
	// Email - the users email
	// required: true
	// example: string
	Email string `json:"email"`
}

// UserResponse - a user object that we can send as json
type UserResponse struct {
	// Email - user's Email
	// example: string
	Email string `json:"email"`
}

const registerUserQuery = `INSERT INTO membership.users(
	email, password)
	VALUES ($1, $2);`
const getUserPasswordQuery = `SELECT password from membership.users where email=$1`
const getUserQuery = `SELECT email from membership.users where email=$1`

// RegisterUser register a user in the db
func (db *Database) RegisterUser(email string, password string) error {
	if len(password) == 0 {
		return fmt.Errorf("not a valid password")
	}

	if len(email) == 0 {
		return fmt.Errorf("not a valid email")
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	// Next, insert the email, along with the hashed password into the database
	rows, err := db.pool.Query(context.Background(), registerUserQuery, email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("conn.Query failed: %s", err)
	}

	defer rows.Close()

	return nil
}

// UserSignin - user login
func (db *Database) UserSignin(email string, password string) error {
	// We create another instance of `Credentials` to store the credentials we get from the database
	storedCreds := &Credentials{}

	// Get the existing entry present in the database for the given user
	row := db.pool.QueryRow(context.Background(), getUserPasswordQuery, email).Scan(&storedCreds.Password)
	if row == pgx.ErrNoRows {
		return fmt.Errorf("Unauthorized")
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(password)); err != nil {
		// If the two passwords don't match, return a 401 status
		return fmt.Errorf("Unauthorized: %s", err)
	}

	return nil
}

// GetUser returns the currently logged in user
func (db *Database) GetUser(email string) (UserResponse, error) {
	var userResponse UserResponse

	row := db.pool.QueryRow(context.Background(), getUserQuery, email).Scan(&userResponse.Email)
	if row == pgx.ErrNoRows {
		return userResponse, fmt.Errorf("error getting user")
	}
	return userResponse, nil
}

package dbstore

var userDbMethod UserDatabaseMethod

// UserDatabaseMethod -- method container that holds the extension methods to query the user table
type UserDatabaseMethod struct{}

func (user *UserDatabaseMethod) registerUser() string {
	const registerUserQuery = `INSERT INTO membership.users(
		email, password)
		VALUES ($1, $2);`

	return registerUserQuery
}

func (user *UserDatabaseMethod) getUserPassword() string {
	const getUserPasswordQuery = `SELECT password from membership.users where email=$1`

	return getUserPasswordQuery
}

func (user *UserDatabaseMethod) getUser() string {
	const getUserQuery = `SELECT email from membership.users where email=$1`

	return getUserQuery
}
